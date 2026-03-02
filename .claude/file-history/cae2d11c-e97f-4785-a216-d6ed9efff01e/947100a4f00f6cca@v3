// ============================================
// Distributed Bank - Rust Implementation
// ============================================
// 고성능 분산 은행 시스템

use dashmap::DashMap;
use std::sync::Arc;
use tokio::sync::RwLock;
use uuid::Uuid;
use tracing::info;

#[derive(Clone, Debug)]
pub struct Account {
    pub id: String,
    pub owner: String,
    pub balance: i64, // 센트 단위
    pub created_at: u64,
}

#[derive(Clone, Debug)]
pub struct Transaction {
    pub tx_id: String,
    pub from: String,
    pub to: String,
    pub amount: i64,
    pub timestamp: u64,
    pub status: String,
}

pub struct BankSystem {
    accounts: Arc<DashMap<String, Account>>,
    transactions: Arc<RwLock<Vec<Transaction>>>,
    total_transactions: Arc<std::sync::atomic::AtomicUsize>,
}

impl BankSystem {
    /// 새 은행 시스템 생성
    pub fn new() -> Self {
        info!("🏦 Bank System initialized");

        BankSystem {
            accounts: Arc::new(DashMap::new()),
            transactions: Arc::new(RwLock::new(Vec::new())),
            total_transactions: Arc::new(std::sync::atomic::AtomicUsize::new(0)),
        }
    }

    /// 계좌 생성
    pub async fn create_account(&self, account_id: &str, initial_balance: i64) {
        let created_at = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .map(|d| d.as_secs())
            .unwrap_or(0);  // 시간 생성 실패 시 0 사용 (거의 발생하지 않음)

        let account = Account {
            id: account_id.to_string(),
            owner: format!("Owner of {}", account_id),
            balance: initial_balance,
            created_at,
        };

        self.accounts.insert(account_id.to_string(), account);
        info!("✅ Account created: {}", account_id);
    }

    /// 계좌 조회
    pub async fn get_account(&self, account_id: &str) -> Option<Account> {
        self.accounts.get(account_id).map(|r| r.clone())
    }

    /// 송금
    pub async fn transfer(
        &self,
        from: &str,
        to: &str,
        amount: i64,
    ) -> Result<String, String> {
        // 수수료 (50센트)
        let fee = 50i64;
        let total = amount + fee;

        // 발신자 잔액 확인
        let from_account = self
            .accounts
            .get_mut(from)
            .ok_or_else(|| "Sender account not found".to_string())?;

        if from_account.balance < total {
            return Err("Insufficient balance".to_string());
        }

        from_account.balance -= total;

        // 수신자 잔액 업데이트
        if let Some(mut to_account) = self.accounts.get_mut(to) {
            to_account.balance += amount;
        } else {
            return Err("Recipient account not found".to_string());
        }

        // 거래 기록
        let tx_id = Uuid::new_v4().to_string();

        let timestamp = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .map(|d| d.as_secs())
            .unwrap_or(0);  // 시간 생성 실패 시 0 사용

        let transaction = Transaction {
            tx_id: tx_id.clone(),
            from: from.to_string(),
            to: to.to_string(),
            amount,
            timestamp,
            status: "completed".to_string(),
        };

        {
            let mut txs = self.transactions.write().await;
            txs.push(transaction);
        }

        self.total_transactions
            .fetch_add(1, std::sync::atomic::Ordering::Relaxed);

        info!(
            "✅ Transfer: {} -> {} ({} cents)",
            from, to, amount
        );

        Ok(tx_id)
    }

    /// 입금
    pub async fn deposit(&self, account_id: &str, amount: i64) -> Result<(), String> {
        if let Some(mut account) = self.accounts.get_mut(account_id) {
            account.balance += amount;
            info!("✅ Deposit: {} ({} cents)", account_id, amount);
            Ok(())
        } else {
            Err("Account not found".to_string())
        }
    }

    /// 출금
    pub async fn withdraw(&self, account_id: &str, amount: i64) -> Result<(), String> {
        if let Some(mut account) = self.accounts.get_mut(account_id) {
            if account.balance < amount {
                return Err("Insufficient balance".to_string());
            }

            account.balance -= amount;
            info!("✅ Withdrawal: {} ({} cents)", account_id, amount);
            Ok(())
        } else {
            Err("Account not found".to_string())
        }
    }

    /// 거래 조회
    pub async fn get_transactions(&self) -> Vec<Transaction> {
        self.transactions.read().await.clone()
    }

    /// 통계
    pub fn get_stats(&self) -> BankStats {
        let total_accounts = self.accounts.len();
        let total_transactions = self
            .total_transactions
            .load(std::sync::atomic::Ordering::Relaxed);

        let total_balance: i64 = self
            .accounts
            .iter()
            .map(|entry| entry.value().balance)
            .sum();

        BankStats {
            total_accounts,
            total_transactions,
            total_balance,
        }
    }

    /// 상태 보고
    pub fn print_status(&self) {
        let stats = self.get_stats();

        info!("\n💰 Bank Statistics:");
        info!("├─ Total Accounts: {}", stats.total_accounts);
        info!("├─ Total Transactions: {}", stats.total_transactions);
        info!("└─ Total Balance: ${:.2}", stats.total_balance as f64 / 100.0);
    }
}

impl Clone for BankSystem {
    fn clone(&self) -> Self {
        BankSystem {
            accounts: self.accounts.clone(),
            transactions: self.transactions.clone(),
            total_transactions: self.total_transactions.clone(),
        }
    }
}

#[derive(Debug, Clone)]
pub struct BankStats {
    pub total_accounts: usize,
    pub total_transactions: usize,
    pub total_balance: i64,
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_create_account() {
        let bank = BankSystem::new();
        bank.create_account("test_account", 100_000_00).await;

        let account = bank.get_account("test_account").await;
        assert!(account.is_some());
        assert_eq!(account.unwrap().balance, 100_000_00);
    }

    #[tokio::test]
    async fn test_transfer() {
        let bank = BankSystem::new();
        bank.create_account("alice", 100_000_00).await;
        bank.create_account("bob", 50_000_00).await;

        let result = bank.transfer("alice", "bob", 30_000_00).await;
        assert!(result.is_ok());

        let alice = bank.get_account("alice").await.unwrap();
        let bob = bank.get_account("bob").await.unwrap();

        assert_eq!(alice.balance, 69_999_50); // 100000 - 30000 - 50
        assert_eq!(bob.balance, 80_000_00); // 50000 + 30000
    }

    #[tokio::test]
    async fn test_deposit() {
        let bank = BankSystem::new();
        bank.create_account("test", 0).await;

        let result = bank.deposit("test", 50_000_00).await;
        assert!(result.is_ok());

        let account = bank.get_account("test").await.unwrap();
        assert_eq!(account.balance, 50_000_00);
    }

    #[tokio::test]
    async fn test_withdraw() {
        let bank = BankSystem::new();
        bank.create_account("test", 100_000_00).await;

        let result = bank.withdraw("test", 30_000_00).await;
        assert!(result.is_ok());

        let account = bank.get_account("test").await.unwrap();
        assert_eq!(account.balance, 70_000_00);
    }

    #[tokio::test]
    async fn test_insufficient_balance() {
        let bank = BankSystem::new();
        bank.create_account("test", 10_000_00).await;

        let result = bank.withdraw("test", 50_000_00).await;
        assert!(result.is_err());
    }
}
