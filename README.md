# score-keep-api
Golang API for handling `score-keep` clients

## Local development

For the `calendar` service you will need Google [service account credentials](https://cloud.google.com/iam/docs/service-account-creds) saved at the root of the directory under `gcp-creds.json`. Reference the `gcp-creds.example.json` file for an example.

```
make local
```

Will use `docker-compose` to start the frontend and backend.

```
make test coverage
```
Will run the tests for every service (excluding web) and output their respective coverage percentages.

## Deployment

This app is hosted on GCP virtual machines and is built using `terraform`. 

There are dependencies in order to properly build the nomad cluster. First we'll need to build the image needed to host the app. 
First install [packer](https://developer.hashicorp.com/packer/tutorials/docker-get-started/get-started-install-cli) for creating the proper image for the VM.

```
cd deploy
packer build -var-file=variables.hcl image.pkr.hcl
// save output for the variables.hcl file
// example output machine image name: hashistack-20240503172529

mv variables.hcl.example variables.hcl
```
Next you will need to update the fields in the newly created `variables.hcl` file.
Once updated apply the `terraform`.
```
terraform apply -var-file=variables.hcl 
```
After success, run `./post-setup.sh` that will handle configuring the ACL tokens correctly for nomad nad consul.

Once the cluster is successfully created you can run the individual jobs to deploy their respective services.
```
nomad job run ./deploy/jobs/score-keep-web-api.hcl
```

