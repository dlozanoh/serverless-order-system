# 🛒 Serverless Order Processing System

A fully serverless order management backend built using **Go**, **AWS Lambda**, **API Gateway**, **SQS**, **DynamoDB**, and **S3**, deployed via the **AWS SAM** framework.

## 📦 Features

- **POST /orders** – Receive new orders and enqueue to SQS
- **SQS → Lambda** – Process order, store in DynamoDB, generate PDF receipt
- **S3** – Upload PDF receipts
- **GET /receipt/{orderId}** – Retrieve signed URL to download receipt (expires in 5 minutes)
- Fully serverless and scalable
- Written entirely in **Go**

---

## 🛠️ Tech Stack

| Component      | Technology          |
|----------------|---------------------|
| Infrastructure | AWS SAM (Serverless Application Model) |
| Language       | Go (1.20+)           |
| API Gateway    | REST API             |
| Queue          | Amazon SQS           |
| Persistence    | Amazon DynamoDB      |
| Storage        | Amazon S3 (receipts) |
| Receipt Format | PDF (via `gofpdf`)   |

---

## 📁 Project Structure

```text
serverless-order-system/
├── template.yaml       # SAM template (infrastructure)
├── Makefile            # Build & deploy commands
├── cmd/                # Go Lambda source code
│ ├── receive_order/    # API Gateway → SQS
│ ├── process_order/    # SQS → DynamoDB & S3
│ └── get_receipt/      # API Gateway → signed URL
├── events/             # Sample test events
├── README.md
```

---

🔌 API Endpoints
POST /orders
Submit a new order.

Request Body:

{
  "orderId": "12345",
  "customerName": "Alice"
}

GET /receipt/{orderId}
Get a signed PDF receipt URL for the given order (valid for 5 minutes).

Response:

{
  "signedReceiptUrl": "https://bucket-name.s3.amazonaws.com/…"
}


📸 Receipt Sample
A clean PDF receipt is generated with:

Order ID

Customer Name

Timestamp

Thank-you note

Powered by gofpdf.

🧪 Local Testing
Test Lambda functions locally using SAM CLI:


sam local invoke ReceiveOrderFunction --event events/receive_order_event.json
sam local invoke ProcessOrderFunction --event events/process_order_event.json
sam local invoke GetReceiptFunction --event events/get_receipt_event.json


✅ TODO (Optional Enhancements)
 Add email notifications via SES

 Frontend order dashboard (React/Vite)

 Add unit & integration tests

 Add order validation schema

📄 License
MIT © YourName

🙋‍♂️ Support
Feel free to open an issue or reach out for improvements, bugs, or questions.

yaml
Copiar
Editar

---

Would you like me to include this README in your project and provide a downloadable `.zip`? Or push it to a GitHub repo?








Preguntar a ChatGPT



