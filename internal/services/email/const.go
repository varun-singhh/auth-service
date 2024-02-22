package email

// Common email template
const customEmailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.header}}</title>
<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
    }
    .container {
        max-width: 600px;
        margin: 20px auto;
        background-color: #fff;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }
    h2 {
        color: #333;
    }
    p {
        color: #666;
        line-height: 1.6;
    }
    .button {
        display: inline-block;
        background-color: #007bff;
        color: #fff;
        text-decoration: none;
        padding: 10px 20px;
        border-radius: 5px;
        margin-top: 20px;
    }
    .button:hover {
        background-color: #0056b3;
    }
</style>
</head>
<body>
<div class="container">
    <h2>{{.subHeader}}</h2>
    <p>{{.title}}</p>
    <p>{{.subTitle}}</p>
    <a href="{{.buttonLink}}" class="button">{{.buttonLabel}}</a>
    <p>{{.body}}</p>
    <p>{{.footer}}</p>
</div>
</body>
</html>
`

const verifyCodeEmailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Email Verification Code</title>
<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
    }
    .container {
        margin: 20px auto;
        background-color: #fff;
        padding: 20px;
        border-radius: 10px;
    }
    h2 {
        color: #333;
    }
    p {
        color: #666;
        line-height: 1.6;
    }
    .verification-code {
        font-size: 24px;
        font-weight: bold;
        padding: 10px 20px;
        color: #808080;
        border-radius: 5px;
        margin-top: 20px;
    }
</style>
</head>
<body>
<div class="container">
    <h2>Email Verification Code</h2>
    <p>Dear {{.user}},</p>
    <p>Your verification code is:</p>
    <div class="verification-code">{{.code}}</div>
        <p class="expires">This code expires in 2 minutes. Please use it promptly.</p>
    <p>Please use this code to verify your email address.</p>
    <p>If you didn't request this verification code, you can safely ignore this email.</p>
    <p>Thank you,<br>AIMSS Chamiana Shimla Himachal Pradesh<br>Phone: 01773501627,01773501628<br>Website: <a href="http://www.aimsschamiana.edu.in/">aimsschamiana.edu.in</a></p>
    
</div>
</body>
</html>
`

const resetPasswordEmailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Password Reset</title>
<style>
    body {
        font-family: Arial, sans-serif;
        background-color: #f4f4f4;
        margin: 0;
        padding: 0;
    }
    .container {
        max-width: 600px;
        margin: 20px auto;
        background-color: #fff;
        padding: 20px;
        border-radius: 10px;
        box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    }
    h2 {
        color: #333;
    }
    p {
        color: #666;
        line-height: 1.6;
    }
    .button {
        display: inline-block;
        background-color: #007bff;
        color: #fff;
        text-decoration: none;
        padding: 10px 20px;
        border-radius: 5px;
        margin-top: 20px;
    }
    .button:hover {
        background-color: #0056b3;
    }
</style>
</head>
<body>
<div class="container">
    <h2>Password Reset</h2>
    <p>Dear {{.user}},</p>
    <p>We received a request to reset your password. If you did not make this request, you can ignore this email.</p>
    <p>To reset your password, click the button below:</p>
    <a href={{.resetLink}} class="button">Reset Password</a>
    <p>If the button above doesn't work, you can also copy and paste the following link into your browser:</p>
    <p>{{.resetLink}}</p>
    <p>This link will expire in 24 hours for security reasons.</p>
    <p>If you have any questions, please contact our support team.</p>
    <p>Thank you,<br>AIMSS Chamiana Shimla Himachal Pradesh<br>Phone: 01773501627,01773501628<br>Website: <a href="http://www.aimsschamiana.edu.in/">aimsschamiana.edu.in</a></p>
</div>
</body>
</html>
`
