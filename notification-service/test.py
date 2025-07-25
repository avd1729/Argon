import boto3
from dotenv import load_dotenv
import os

load_dotenv()
email_id = os.getenv("EMAIL_ID")

ses = boto3.client("ses", region_name="us-east-1")  

response = ses.send_email(
    Source=email_id,  
    Destination={
        "ToAddresses": [email_id],
    },
    Message={
        "Subject": {"Data": "Hello from SES"},
        "Body": {"Text": {"Data": "This is a test email from SES using Python!"}},
    }
)

print("Email sent! Message ID:", response['MessageId'])
