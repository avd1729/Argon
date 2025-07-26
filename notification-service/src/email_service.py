import boto3
from dotenv import load_dotenv
import os

load_dotenv()
email_id = os.getenv("EMAIL_ID")
ses = boto3.client("ses", region_name="us-east-1")  

import boto3
from dotenv import load_dotenv
import os

load_dotenv()
email_id = os.getenv("EMAIL_ID")
ses = boto3.client("ses", region_name="us-east-1")  

def send_email(body: str):
    response = ses.send_email(
        Source=email_id,
        Destination={
            "ToAddresses": [email_id],
        },
        Message={
            "Subject": {
                "Data": "CI Run Results",
                "Charset": "UTF-8"
            },
            "Body": {
                "Text": {
                    "Data": body,
                    "Charset": "UTF-8"
                }
            }
        }
    )
    print("Email sent! Message ID:", response['MessageId'])
