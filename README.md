# CryptoCurrencyMailer
The Crypto Currency Mailer project is a tool that allows users to stay updated with the latest exchange rates of cryptocurrencies, specifically Bitcoin (BTC) to Ukrainian Hryvnia (UAH). 
It utilizes BinanceAPI and CoinAPI to fetch real-time exchange rate data and provides a subscription feature for users to receive email notifications when the exchange rate updates.

Key Features:

* **Real-time Exchange Rate**: The project integrates with Binance and Coin APIs to fetch the current BTC to UAH exchange rate. It ensures that users receive up-to-date information on the cryptocurrency market.

* **Subscription Management**: Users can subscribe to the service by providing their email addresses. The project uses a .txt file to store the list of subscribers, making it easy to manage and update the subscriber list.

* **Email Notifications**: Subscribers receive email notifications whenever there is a change in the BTC to UAH exchange rate. The project leverages SMTP (Simple Mail Transfer Protocol) to send personalized emails containing the updated exchange rate information.

# Installation

### Gmail SMTP Server
The Crypto Currency Mailer project employs the SMTP protocol for sending emails. To use this protocol, it is necessary to authenticate using your email credentials and an application-specific password. You can refer to [Guide Link](https://kinsta.com/blog/gmail-smtp-server/) for instructions on generating an app password.

### Update Environment Variables
Inside the **.env** file there is a placeholder of values, including SENDER_EMAIL and SENDER_PASSWORD variables. This ensures that the mailer has the necessary authentication credentials to establish a secure connection with the SMTP server.
Also, there is a COIN_API_KEY variable that contain unique key, that give you access to make requests to CoinAPI to retrieve rates.

Due to the security, I removed these values in advance, so you will need to paste your own there.

# Docker
To deploy the project I use docker compose, so I recommend the same to you. If you haven't installed it yet, there is a link: [Install Docker](https://docs.docker.com/desktop/).
```
docker compose up --build
```

# Third-party APIs
The Binance API provides a comprehensive set of endpoints and documentation that allows developers to access various cryptocurrency-related data and services.
These endpoints specifically designed for fetching market data, including ticker information for different trading pairs. The application uses the "api.binance.com/api/v3/ticker/price?symbol=BTCUAH" endpoint to retrieve the price of the BTC to UAH trading pair, but it also can work with another currency pairs.

To ensure a fail-safe mechanism, the application can integrate CoinAPI as an alternative rate provider.
The application uses "rest.coinapi.io/v1/exchangerate/BTC/UAH" endpoint to retrieve BTC-UAH exchange rate.
For example, if primary provider, such as BinanceAPI, encounters issues or fails to provide the required exchange rate, the application will automatically pass the request to next provider - CoinAPI. 
Overall, the integration of Chain of Responsibility pattern provides a safety net for the application.  

# API Endpoints

### Rate
* **Get** localhost:3000/api/rate
  
    Get current cryptocurrency exchange rate

    Parameters

    Responses
    * 200 - returns current exchange rate
    * 400 - Bad Request

### Mailer

* **Post** localhost:3000/api/sendEmails

    Parameters

    Responses
    * 200 - returns current exchange rate
    * 400 - Bad Request
    * 500 - Internal Server Error


* **Post** localhost:3000/api/subscribe

    Parameters

    Body (email string - email, that you want to subscribe)

    Responses
    * 200 - returns current exchange rate
    * 400 - Bad Request
    * 409 - Conflict
    * 500 - Internal Server Error

# Tests
To run tests use this command: go test ./...