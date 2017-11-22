"""
An Airflow pipeline for running our E2E tests.
"""

# TODO(jlewi): We should split setup into two steps; create cluster and setup
# cluster. The cluster can be created in parallel with building the artifacts
# which should speed things up.

from datetime import datetime
import logging
import os
import uuid

from airflow import DAG
from airflow.operators import PythonOperator
from py import util
from google.cloud import storage  # pylint: disable=no-name-in-module
import six
import yaml

default_args = {
  'owner': 'airflow',
    'depends_on_past': False,
    'start_date': datetime(2015, 6, 1),
    'email': ['airflow@airflow.com'],
    'email_on_failure': False,
    'email_on_retry': False,
    'retries': 1,
}

dag = DAG(
  # Set schedule_interval to None
  'tf_k8s_tests', default_args=default_args,
  # TODO(jlewi): Should we schedule a regular run? Right now its
  # manually triggered by PROW.
  schedule_interval=None)

# Default name for the repo organization and name.
# This should match the values used in Go imports.
DEFAULT_REPO_OWNER = "tensorflow"
DEFAULT_REPO_NAME = "k8s"

# GCS path to use for each dag run
GCS_RUNS_PATH = "gs://mlkube-testing-airflow/runs"

GCB_PROJECT = "mlkube-testing"

def run_path(dag_id, run_id):
  return os.path.join(GCS_RUNS_PATH, dag_id.replace(":", "_"),
                      run_id.replace(":", "_"))

class FakeDagrun(object):
  def __init__(self):
    self.dag_id = "tf_k8s_tests"
    self.run_id = "test_run"
    self.conf = {}

def build_images(dag_run=None, ti=None, **_kwargs):
  """
  Args:
    dag_run: A DagRun object. This is passed in as a result of setting
      provide_context to true for the operator.
  """
  # Create a temporary directory suitable for checking out and building the
  # code.
  if not dag_run:
    # When running via airflow test dag_run isn't set
    logging.warn("Using fake dag_run")
    dag_run = FakeDagrun()

  logging.info("dag_id: %s", dag_run.dag_id)
  logging.info("run_id: %s", dag_run.run_id)

  gcs_path = run_path(dag_run.dag_id, dag_run.run_id)
  logging.info("gcs_path %s", gcs_path)

  conf = dag_run.conf
  if not conf:
    conf = {}
  logging.info("conf=%s", conf)
  artifacts_path = conf.get("ARTIFACTS_PATH", gcs_path)
  logging.info("artifacts_path %s", artifacts_path)

  # Make sure pull_number is a string
  pull_number = "{0}".format(conf.get("PULL_NUMBER", ""))
  args = ["python", "-m", "py.release"]
  if pull_number:
    commit = conf.get("PULL_PULL_SHA", "")
    args.append("pr")
    args.append("--pr=" + pull_number)
    if commit:
      args.append("--commit=" + commit)
  else:
    commit = conf.get("PULL_BASE_SHA", "")
    args.append("postsubmit")
    if commit:
      args.append("--commit=" + commit)

  dryrun = bool(conf.get("dryrun", False))

  build_info_file = os.path.join(gcs_path, "build_info.yaml")
  args.append("--build_info_path=" + build_info_file)
  args.append("--releases_path=" + gcs_path)
  args.append("--project=" + GCB_PROJECT)

  # We want subprocess output to bypass logging module otherwise multiline
  # output is squashed together.
  util.run(args, use_print=True, dryrun=dryrun)

  # Read the output yaml and publish relevant values to xcom.
  if not dryrun:
    gcs_client = storage.Client(project=GCB_PROJECT)
    logging.info("Reading %s", build_info_file)
    bucket_name, build_path = util.split_gcs_uri(build_info_file)
    bucket = gcs_client.get_bucket(bucket_name)
    blob = bucket.blob(build_path)
    contents = blob.download_as_string()
    build_info = yaml.load(contents)
  else:
    build_info = {
      "image": "gcr.io/dryrun/dryrun:latest",
      "commit": "1234abcd",
      "helm_chart": "gs://dryrun/dryrun.latest.",
    }
  for k, v in six.iteritems(build_info):
    logging.info("xcom push: %s=%s", k, v)
    ti.xcom_push(key=k, value=v)

def setup_cluster(dag_run=None, ti=None, **_kwargs):
  conf = dag_run.conf
  if not conf:
    conf = {}

  dryrun = bool(conf.get("dryrun", False))

  chart = ti.xcom_pull("build_images", key="helm_chart")

  now = datetime.now()
  cluster = "e2e-" + now.strftime("%m%d-%H%M-") + uuid.uuid4().hex[0:4]

  logging.info("conf=%s", conf)
  artifacts_path = conf.get("ARTIFACTS_PATH",
                            run_path(dag_run.dag_id, dag_run.run_id))
  logging.info("artifacts_path %s", artifacts_path)

  # Gubernator only recognizes XML files whos name matches
  # junit_[^_]*.xml which is why its "setupcluster" and not "setup_cluster"
  junit_path = os.path.join(artifacts_path, "junit_setupcluster.xml")
  logging.info("junit_path %s", junit_path)

  args = ["python", "-m", "py.deploy", "setup"]
  args.append("--cluster=" + cluster)
  args.append("--junit_path=" + junit_path)
  args.append("--project=" + GCB_PROJECT)
  args.append("--chart=" + chart)

  # We want subprocess output to bypass logging module otherwise multiline
  # output is squashed together.
  util.run(args, use_print=True, dryrun=dryrun)

  values = {
    "cluster": cluster,
  }
  for k, v in six.iteritems(values):
    logging.info("xcom push: %s=%s", k, v)
    ti.xcom_push(key=k, value=v)

def run_tests(dag_run=None, ti=None, **_kwargs):
  conf = dag_run.conf
  if not conf:
    conf = {}

  dryrun = bool(conf.get("dryrun", False))

  cluster = ti.xcom_pull("setup_cluster", key="cluster")

  logging.info("conf=%s", conf)
  artifacts_path = conf.get("ARTIFACTS_PATH",
                            run_path(dag_run.dag_id, dag_run.run_id))
  logging.info("artifacts_path %s", artifacts_path)
  junit_path = os.path.join(artifacts_path, "junit_e2e.xml")
  logging.info("junit_path %s", junit_path)
  ti.xcom_push(key="cluster", value=cluster)

  args = ["python", "-m", "py.deploy", "test"]
  args.append("--cluster=" + cluster)
  args.append("--junit_path=" + junit_path)
  args.append("--project=" + GCB_PROJECT)

  # We want subprocess output to bypass logging module otherwise multiline
  # output is squashed together.
  util.run(args, use_print=True, dryrun=dryrun)

def teardown_cluster(dag_run=None, ti=None, **_kwargs):
  conf = dag_run.conf
  if not conf:
    conf = {}

  dryrun = bool(conf.get("dryrun", False))

  cluster = ti.xcom_pull("setup_cluster", key="cluster")

  gcs_path = run_path(dag_run.dag_id, dag_run.run_id)

  artifacts_path = conf.get("ARTIFACTS_PATH", gcs_path)
  logging.info("artifacts_path %s", artifacts_path)

  junit_path = os.path.join(artifacts_path, "junit_teardown.xml")
  logging.info("junit_path %s", junit_path)
  ti.xcom_push(key="cluster", value=cluster)

  args = ["python", "-m", "py.deploy", "teardown"]
  args.append("--cluster=" + cluster)
  args.append("--junit_path=" + junit_path)
  args.append("--project=" + GCB_PROJECT)

  # We want subprocess output to bypass logging module otherwise multiline
  # output is squashed together.
  util.run(args, use_print=True, dryrun=dryrun)

build_op = PythonOperator(
  task_id='build_images',
    provide_context=True,
    python_callable=build_images,
    dag=dag)

setup_cluster_op = PythonOperator(
  task_id='setup_cluster',
    provide_context=True,
    python_callable=setup_cluster,
    dag=dag)

setup_cluster_op.set_upstream(build_op)

run_tests_op = PythonOperator(
  task_id='run_tests',
    provide_context=True,
    python_callable=run_tests,
    dag=dag)

run_tests_op.set_upstream(setup_cluster_op)

teardown_cluster_op = PythonOperator(
  task_id='teardown_cluster',
    provide_context=True,
    python_callable=teardown_cluster,
    dag=dag)

teardown_cluster_op.set_upstream(run_tests_op)
