# Stori Transaction Reports Service

## Getting Started

This is a service for processing transaction files and sending out monthly summary reports. It's written in Go and can be deployed on **ECS** or in **Lambda**. Below you'll find everything you need to get it running locally or on AWS. This guide will help you get started, but assumes some basic familiarity with Docker, AWS, and Go.

## Clone the Code

To get the code, just clone the repo and cd into it:

```sh
$ git clone 
$ cd StoriTransactionReports
```

## Database Setup

This service needs a PostgreSQL database to store transactions. Here's the minimal table definition to make it work:

```sql
CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  date VARCHAR(5),
  amount NUMERIC
);
```
- **`id`**: Unique transaction identifier.
- **`date`**: Date of the transaction, formatted as `month/day` (e.g., `7/15`).
- **`amount`**: Transaction amount, either positive or negative (e.g., `+60.5`, `-10.3`).

## Running Locally with Docker

To run the service locally using Docker, you'll need to pass in some environment variables. Here's how to do it. Just use the command below and make sure to set the correct values for your environment:

```sh
$ docker build -t stori-reports .
$ docker run -p 8080:8080 \
  -e DB_HOST=<your-db-host> \
  -e DB_PORT=<your-db-port> \
  -e DB_USER=<your-db-user> \
  -e DB_PASSWORD=<your-db-password> \
  -e DB_NAME=<your-db-name> \
  -e SENDER_EMAIL=<your-smtp-email> \
  -e SENDER_PASSWORD=<your-smtp-password> \
  -e SMTP_HOST=<your-smtp-host> \
  -e AWS_REGION=<your-aws-region> \
  stori-reports
```

## Deploying to AWS

There are two ways to deploy the service to **AWS**: **ECS** or **Lambda**:

### Deploying to ECS
1. **Build and Push to ECR**:
   ```sh
   aws ecr create-repository --repository-name stori-reports
   docker tag stori-reports:latest <your-account-id>.dkr.ecr.<region>.amazonaws.com/stori-reports:latest
   docker push <your-account-id>.dkr.ecr.<region>.amazonaws.com/stori-reports:latest
   ```

2. **ECS Task Definition**:
    - Create a task definition that points to the **ECR image**.
    - Add the **environment variables** to the task definition (same ones listed above).

3. **IAM Role**:
    - If you want to be able to see the container logs in **Cloud Watch**, make sure you have a task **execution role** that has permissions to access it.

4. **Run It**:
    - Create an ECS service using the task definition. Make sure the **security groups** and **networking** are set up correctly to open up the port 8080 to the container and let the ECS tasks connect to your RDS instance(If the DB is set up on RDS).

5. **Logs**:
    - Use **CloudWatch Logs** to monitor the containers.



### Deploying to Lambda
1. **Build and package the executable**: 


   **Windows Powershell**
   ```sh
   docker run --rm -v ${PWD}:/usr/src/myapp -w /usr/src/myapp golang:latest /bin/sh -c "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/lambda/lambda_main.go"
   Compress-Archive -Path ./bootstrap,./assets -DestinationPath lambda_package.zip
   ```

   **Linux**
   ```sh
   go build -o bootstrap ./cmd/lambda/lambda_main.go
   zip lambda_package.zip bootstrap assets
   ```

2. **Create and set up Lambda Function**:
   - Create a new lambda function, select **Author from scratch**.
   - Select **Amazon Linux 2023** as the runtime.
   - Choose the correct architecture, given the executable file you created(For instance, if you used the Windows Powershell command in the previous step, it'd be **x86_64**)
   - Add the **environment variables** to the lambda function (same ones listed above).
   - In **Code source**, choose **Upload from .zip file**, and upload the zip file created in the previous step.
   - In **Runtime settings**, use **bootstrap** as the **Handler**. 

3. **Create and configure S3 bucket**:
   - Create a new S3 bucket.
   - Give the appropriate **permissions** to the bucket(For instance, if no one else is going to be uploading files to it, set up the bucket policy so only your IP can access it).
   - Add a new **Event notification** that triggers the lambda created in the previous step whenever a new **csv** file is uploaded to the bucket.

4. **Logs**:
   - Use **CloudWatch Logs** to monitor the containers.


# Using the service

Given this is a proof of concept, we wanted to consider different possible ways of triggering the service, before deciding on which one to keep, that's why it works slightly differently when deployed on **ECS** and when deployed on **Lambda**.

When Deployed on **ECS**, the service receives the CSV file in the body of the request, and the recipient in a header. When deployed on **Lambda**, the service is triggered automatically when a new file is uploaded to a specific **S3** bucket, and the name of the csv file is used to choose a recipient.

## Using the service if running locally on docker or on ECS

When sending requests to the service, just send **POST** requests to the port 8080 with the following headers:

**Content-Type** : **text/csv**

**recipient** : **< Email report recipient >**

And the body containing the **raw CSV**, for example:

```
ID,Date,Amount
5,7/15,+60.5
6,7/28,-10.3
7,8/2,-20.46
9,8/13,+10
```

## Using the service if running Lambda

Simply upload a csv file to the **S3 bucket** we created and set up when deploying the service, the Lambda will be then triggered automatically and email the address in the name of the csv file. 

For instance, if we upload **adavilam@unal.edu.co.csv** to the S3 bucket, an email will be sent to **adavilam@unal.edu.co**.