# Concurrent Wallet Order System

A robust Go-based financial trading platform that manages user wallets, stock trading, and investment portfolios with concurrent-safe operations using MongoDB as the database backend.

## Project Overview

This system provides a complete trading platform with the following core features:

- **User Management**: User registration and authentication with password hashing
- **Wallet Management**: Deposit and withdraw funds, track transaction history
- **Stock Management**: Create and manage stocks with current pricing
- **Order Processing**: Buy and sell stocks with automatic portfolio updates
- **Portfolio Tracking**: View holdings with real-time valuation

## Architecture

The project follows a layered architecture pattern:

```
cmd/
  └── main.go                 # Application entry point
internal/
  ├── config/                # MongoDB configuration and indexing
  │   ├── indexes.go        # Database index definitions
  │   └── mongo.go          # MongoDB connection setup
  ├── handlers/             # HTTP request handlers (API layer)
  ├── middleware/           # HTTP middleware (auth_middleware.go - empty)
  ├── models/               # Data models/entities
  ├── repo/                 # Data access layer (repositories)
  ├── services/             # Business logic layer
  └── validators/           # Request validation (empty)
```

## Technology Stack

- **Language**: Go 1.25.6
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **Database**: MongoDB (go.mongodb.org/mongo-driver v1.17.9)
- **Cryptography**: golang.org/x/crypto (password hashing with bcrypt)

## Database Schema

### Collections

#### Users
- `_id`: ObjectID (Primary Key)
- `name`: User's display name
- `email`: Email (unique index)
- `password`: Bcrypt hashed password
- `walletbalance`: Current wallet balance (float)
- `createdAt`: Timestamp

#### Stocks
- `_id`: ObjectID (Primary Key)
- `symbol`: Stock ticker symbol (unique index)
- `name`: Company/stock name
- `price`: Current stock price
- `createdAt`: Timestamp

#### Portfolio
- `_id`: ObjectID (Primary Key)
- `userId`: Reference to user (compound unique index: userId + symbol)
- `symbol`: Stock symbol
- `quantity`: Number of shares held

#### Orders
- `_id`: ObjectID (Primary Key)
- `userId`: Reference to user (index)
- `symbol`: Stock symbol
- `type`: "BUY" or "SELL"
- `quantity`: Number of shares
- `price`: Price per share at time of order
- `createdAt`: Timestamp

#### Wallets (Transaction History)
- `_id`: ObjectID (Primary Key)
- `userId`: Reference to user
- `method`: "deposit" or "withdraw"
- `amount`: Transaction amount
- `createdAt`: Timestamp

## API Endpoints

### Authentication & Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register new user |
| POST | `/login` | User login |
| GET | `/users` | Get all users |
| GET | `/users/:userId` | Get user details |

**Register Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Login Request:**
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

### Wallet Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/wallet/deposit` | Deposit funds |
| POST | `/wallet/withdraw` | Withdraw funds |
| GET | `/wallet/balance/:userId` | Get wallet balance |
| GET | `/wallet/history/:userId` | Get transaction history |

**Wallet Request:**
```json
{
  "userId": "507f1f77bcf86cd799439011",
  "amount": 1000.50
}
```

### Stock Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/stocks` | Create new stock |
| GET | `/stocks` | Get all stocks |
| GET | `/stocks/:symbol` | Get stock by symbol |

**Create Stock Request:**
```json
{
  "symbol": "AAPL",
  "name": "Apple Inc.",
  "price": 150.75
}
```

### Order Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders/buy` | Place buy order |
| POST | `/orders/sell` | Place sell order |

**Order Request:**
```json
{
  "userId": "507f1f77bcf86cd799439011",
  "symbol": "AAPL",
  "quantity": 10
}
```

### Portfolio

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/portfolio/:userId` | Get user portfolio with valuation |

**Portfolio Response:**
```json
{
  "userId": "507f1f77bcf86cd799439011",
  "holdings": [
    {
      "symbol": "AAPL",
      "stockName": "Apple Inc.",
      "quantity": 10,
      "currentPrice": 150.75,
      "totalValue": 1507.50
    }
  ],
  "totalPortfolioValue": 1507.50
}
```

## Concurrency & Thread Safety

The system implements mutex-based locking to prevent race conditions:

### Critical Sections Protected:

1. **Wallet Service** (`wallet_service.go`)
   - Deposit and Withdraw operations use `sync.Mutex`
   - Ensures atomic balance updates

2. **Order Service** (`order_service.go`)
   - Buy and Sell operations use `sync.Mutex`
   - Protects stock validation, wallet deduction, and portfolio updates

This prevents concurrent requests from causing:
- Double-spending
- Invalid portfolio states
- Lost transactions

## File Structure Details

### Models (`internal/models/`)
- `user.go`: User entity with wallet balance
- `wallet.go`: WalletTransaction entity for audit trail
- `stock.go`: Stock entity with pricing
- `order.go`: Order entity for trade records
- `portfolio.go`: Portfolio holding entity

### Services (`internal/services/`)
Business logic layer implementing:
- **UserService**: Registration/login with bcrypt password hashing
- **WalletService**: Balance management with mutex protection
- **StockService**: Stock creation and retrieval
- **OrderService**: Buy/sell operations with concurrent safety
- **PortfolioService**: Aggregated portfolio view with current valuations

### Repositories (`internal/repo/`)
Data access layer using MongoDB:
- **UserRepository**: User CRUD and balance updates
- **WalletRepository**: Transaction history recording
- **StockRepository**: Stock CRUD operations
- **OrderRepository**: Order recording
- **PortfolioRepository**: Portfolio upsert/retrieval with aggregation pipelines

### Handlers (`internal/handlers/`)
HTTP request handlers implementing REST endpoints:
- **UserHandler**: `/register`, `/login`, `/users` routes
- **WalletHandler**: Wallet operations
- **StockHandler**: Stock management
- **OrderHandler**: Stock trading
- **PortfolioHandler**: Portfolio retrieval

### Configuration (`internal/config/`)
- **mongo.go**: MongoDB connection initialization
- **indexes.go**: Database index creation for performance optimization

## Getting Started

### Prerequisites

- Go 1.25.6 or later
- MongoDB running on `localhost:27017`

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd concurrent-wallet-order-system
```

2. Install dependencies:
```bash
go mod download
```

3. Start MongoDB (if not already running):
```bash
mongod
```

4. Run the application:
```bash
go run ./cmd/main.go
```

The server will start on `http://localhost:8080`

## Key Features

### Security
- Bcrypt password hashing with salting
- Unique email constraints in database
- No password exposure in API responses

### Performance
- MongoDB indexes on frequently queried fields
- Unique compound indexes for data integrity
- Background index creation

### Data Integrity
- Mutex-protected critical sections
- ACID-like operations for financial transactions
- Upsert patterns for portfolio management

### Error Handling
- Validation for amounts (must be > 0)
- Business logic validation (insufficient balance, stock not found)
- HTTP status codes for different error scenarios

## Database Indexes

Key indexes for performance:
- `users.email` (unique)
- `stocks.symbol` (unique)
- `portfolio.userId` + `portfolio.symbol` (unique compound)
- `orders.userId`

## Transaction Flow Examples

### Buy Order Flow:
1. Validate quantity > 0
2. Lock OrderService mutex
3. Verify stock exists
4. Calculate total cost
5. Withdraw funds from wallet (wallet mutex locked)
6. Update portfolio (upsert)
7. Record order in database
8. Unlock mutex

### Portfolio Valuation:
1. Fetch all user holdings from portfolio
2. For each holding:
   - Get current stock price
   - Calculate position value (quantity × price)
3. Sum all positions for total portfolio value

## Configuration

Default MongoDB connection:
- **URI**: `mongodb://localhost:27017`
- **Database**: `wallet_order_system`
- **Server Port**: `8080`

To modify, edit [cmd/main.go](cmd/main.go):
```go
mongoURI := "mongodb://localhost:27017"
dbName := "wallet_order_system"
```

## Project Structure

```
concurrent-wallet-order-system/
├── go.mod                  # Go module definition
├── go.sum                  # Go dependencies checksums
├── cmd/
│   └── main.go            # Application entry point
└── internal/
    ├── config/
    │   ├── indexes.go
    │   └── mongo.go
    ├── handlers/
    │   ├── order_handler.go
    │   ├── portfolio_handler.go
    │   ├── stock_handler.go
    │   ├── user_handler.go
    │   └── wallet_handler.go
    ├── middleware/
    │   └── auth_middleware.go (empty - ready for implementation)
    ├── models/
    │   ├── order.go
    │   ├── portfolio.go
    │   ├── stock.go
    │   ├── user.go
    │   └── wallet.go
    ├── repo/
    │   ├── order_repo.go
    │   ├── portfolio_repo.go
    │   ├── stock_repo.go
    │   ├── user_repo.go
    │   └── wallet_repo.go
    ├── services/
    │   ├── order_service.go
    │   ├── portfolio_service.go
    │   ├── stock_service.go
    │   ├── user_service.go
    │   └── wallet_service.go
    └── validators/
        └── request_validator.go (empty - ready for implementation)
```

## Future Enhancements

- JWT token-based authentication in auth middleware
- Input validation in validators package
- Rate limiting and request throttling
- WebSocket support for real-time price updates
- Advanced portfolio analytics
- Trading notifications
- Dividend handling
- Margin trading support

## License

This project is provided as-is for educational and development purposes.
