// Script para criar transações recentes (últimos 30 dias)
db = db.getSiblingDB('financial-control');

const user = db.users.findOne({email: "vasconcellos@gmail.com"});
const vasconcellosUserId = user._id;

const accounts = db.accounts.find({user_id: vasconcellosUserId}).toArray();
const categories = db.categories.find({user_id: vasconcellosUserId}).toArray();

const incomeCategories = categories.filter(c => c.type === "income");
const expenseCategories = categories.filter(c => c.type === "expense");

print("📊 Criando transações recentes para: " + user.email);

const transactions = [];

// Criar transações nos últimos 30 dias
const now = new Date();
const startDate = new Date(now.getTime() - (30 * 24 * 60 * 60 * 1000));

for (let i = 0; i < 50; i++) {
    const randomTime = startDate.getTime() + Math.random() * (now.getTime() - startDate.getTime());
    const occurredAt = new Date(randomTime);
    
    const isIncome = Math.random() < 0.3;
    const categoryList = isIncome ? incomeCategories : expenseCategories;
    const categoryIndex = Math.floor(Math.random() * categoryList.length);
    const category = categoryList[categoryIndex];
    
    const accountIndex = Math.floor(Math.random() * accounts.length);
    const account = accounts[accountIndex];
    
    let amount;
    if (isIncome) {
        amount = Math.floor(Math.random() * 19500) + 500;
    } else {
        amount = -(Math.floor(Math.random() * 1490) + 10);
    }
    
    transactions.push({
        user_id: vasconcellosUserId,
        account_id: account._id.toString(),
        category_id: category._id.toString(),
        amount: amount,
        currency: "BRL",
        description: isIncome ? "Receita " + i : "Despesa " + i,
        occurred_at: occurredAt,
        status: "completed",
        notes: "",
        tags: isIncome ? ["receita", "recente"] : ["despesa", "recente"],
        created_at: occurredAt,
        updated_at: occurredAt,
        external_ref: "",
        metadata: {}
    });
}

const result = db.transactions.insertMany(transactions);
print("✅ Inseridas " + Object.keys(result.insertedIds).length + " transações recentes do Vasconcellos");
print("Total de transações do Vasconcellos: " + db.transactions.countDocuments({user_id: vasconcellosUserId}));

// Criar transações recentes para o usuário Teste
const testeUser = db.users.findOne({email: "teste@gmail.com"});
if (testeUser) {
    const testeUserId = testeUser._id;
    const testeAccounts = db.accounts.find({user_id: testeUserId}).toArray();
    const testeCategories = db.categories.find({user_id: testeUserId}).toArray();
    
    const testeIncomeCategories = testeCategories.filter(c => c.type === "income");
    const testeExpenseCategories = testeCategories.filter(c => c.type === "expense");
    
    print("\n📊 Criando transações recentes para: " + testeUser.email);
    
    const testeTransactions = [];
    
    for (let i = 0; i < 30; i++) {
        const randomTime = startDate.getTime() + Math.random() * (now.getTime() - startDate.getTime());
        const occurredAt = new Date(randomTime);
        
        const isIncome = Math.random() < 0.4;
        const categoryList = isIncome ? testeIncomeCategories : testeExpenseCategories;
        const categoryIndex = Math.floor(Math.random() * categoryList.length);
        const category = categoryList[categoryIndex];
        
        const accountIndex = Math.floor(Math.random() * testeAccounts.length);
        const account = testeAccounts[accountIndex];
        
        let amount;
        if (isIncome) {
            amount = Math.floor(Math.random() * 15000) + 1000;
        } else {
            amount = -(Math.floor(Math.random() * 1000) + 50);
        }
        
        testeTransactions.push({
            user_id: testeUserId,
            account_id: account._id.toString(),
            category_id: category._id.toString(),
            amount: amount,
            currency: "BRL",
            description: isIncome ? "Receita " + i : "Despesa " + i,
            occurred_at: occurredAt,
            status: "completed",
            notes: "",
            tags: isIncome ? ["receita", "recente"] : ["despesa", "recente"],
            created_at: occurredAt,
            updated_at: occurredAt,
            external_ref: "",
            metadata: {}
        });
    }
    
    const testeResult = db.transactions.insertMany(testeTransactions);
    print("✅ Inseridas " + Object.keys(testeResult.insertedIds).length + " transações recentes do Teste");
    print("Total de transações do Teste: " + db.transactions.countDocuments({user_id: testeUserId}));
} else {
    print("⚠️ Usuário Teste não encontrado. Execute seed_robust_data.js primeiro.");
}

