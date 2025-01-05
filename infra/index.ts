import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

const assumeRole = aws.iam.getPolicyDocument({
  statements: [
    {
      effect: "Allow",
      principals: [
        {
          type: "Service",
          identifiers: ["lambda.amazonaws.com"],
        },
      ],
      actions: ["sts:AssumeRole"],
    },
  ],
});

const role = new aws.iam.Role("lambda_role", {
  name: "lambda_role",
  assumeRolePolicy: assumeRole.then((assumeRole) => assumeRole.json),
});

const lambdaZip = new pulumi.asset.FileArchive("../lambda.zip");

const lambda = new aws.lambda.Function("test_pulumi", {
  runtime: aws.lambda.Runtime.CustomAL2023,
  code: lambdaZip,
  architectures: ["arm64"],
  role: role.arn,
  handler: "bootstrap",
});
