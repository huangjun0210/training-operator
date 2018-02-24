local params = import "../../components/params.libsonnet";
params {
  components+: {
    // Insert component parameter overrides here. Ex:
    // guestbook +: {
    //   name: "guestbook-dev",
    //   replicas: params.global.replicas,
    // },
    workflows+: {
      project: "kubeflow-releasing",
      zone: "us-central1-a",
    },
  },
}
