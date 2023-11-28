# coding: utf-8

"""
    Kubeflow Training SDK

    Python SDK for Kubeflow Training  # noqa: E501

    The version of the OpenAPI document: v1.7.0
    Generated by: https://openapi-generator.tech
"""


from __future__ import absolute_import

import unittest
import datetime

from kubeflow.training.models import *
from kubeflow.training.models.kubeflow_org_v1_paddle_elastic_policy import KubeflowOrgV1PaddleElasticPolicy  # noqa: E501
from kubeflow.training.rest import ApiException

class TestKubeflowOrgV1PaddleElasticPolicy(unittest.TestCase):
    """KubeflowOrgV1PaddleElasticPolicy unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test KubeflowOrgV1PaddleElasticPolicy
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # model = kubeflow.training.models.kubeflow_org_v1_paddle_elastic_policy.KubeflowOrgV1PaddleElasticPolicy()  # noqa: E501
        if include_optional :
            return KubeflowOrgV1PaddleElasticPolicy(
                max_replicas = 56, 
                max_restarts = 56, 
                metrics = [
                    None
                    ], 
                min_replicas = 56
            )
        else :
            return KubeflowOrgV1PaddleElasticPolicy(
        )

    def testKubeflowOrgV1PaddleElasticPolicy(self):
        """Test KubeflowOrgV1PaddleElasticPolicy"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)


if __name__ == '__main__':
    unittest.main()
