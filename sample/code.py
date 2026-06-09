import boto3

ec2 = boto3.client(
    "ec2",
    endpoint_url="http://192.168.1.6:4566",
    region_name="us-east-1",
    aws_access_key_id="test",
    aws_secret_access_key="test"
)

ec2.stop_instances(InstanceIds=["i-5209b79669af77d76"])