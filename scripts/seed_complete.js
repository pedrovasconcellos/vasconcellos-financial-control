// Script completo para criar usuários + todos os dados
db = db.getSiblingDB('financial-control');

function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// ============================================================================
// CRIAR USUÁRIOS
// ============================================================================

print("👥 Verificando/Criando usuários...");

// Verificar se o usuário vasconcellos já existe
let existingVasconcellos = db.users.findOne({email: "vasconcellos@gmail.com"});

let vasconcellosUserId;
if (existingVasconcellos) {
    vasconcellosUserId = existingVasconcellos._id;
    print("✅ Usuário vasconcellos já existe (ID: " + vasconcellosUserId + ")");
} else {
    vasconcellosUserId = generateUUID();
    const vasconcellosUser = db.users.insertOne({
        _id: vasconcellosUserId,
        email: "vasconcellos@gmail.com",
        name: "Vasconcellos",
        default_currency: "BRL",
        cognito_sub: "local-vasconcellos-user",
        created_at: new Date(),
        updated_at: new Date()
    });
    print("✅ Usuário vasconcellos criado (ID: " + vasconcellosUserId + ")");
}

// Verificar se o usuário teste já existe
let existingTeste = db.users.findOne({email: "teste@gmail.com"});

let testeUserId;
if (existingTeste) {
    testeUserId = existingTeste._id;
    print("✅ Usuário teste já existe (ID: " + testeUserId + ")");
} else {
    testeUserId = generateUUID();
    const testeUser = db.users.insertOne({
        _id: testeUserId,
        email: "teste@gmail.com",
        name: "Usuário Teste",
        default_currency: "BRL",
        cognito_sub: "local-teste-user",
        created_at: new Date(),
        updated_at: new Date()
    });
    print("✅ Usuário teste criado (ID: " + testeUserId + ")");
}

// ============================================================================
// EXECUTAR SEED ROBUSTO
// ============================================================================

print("\n📊 Executando seed robusto...");

// Usar os IDs gerados acima
// vasconcellosUserId já foi definido acima

function insertCategory(doc) {
    const id = generateUUID();
    db.categories.insertOne(Object.assign({_id: id}, doc));
    return id;
}

function insertAccount(doc) {
    const id = generateUUID();
    db.accounts.insertOne(Object.assign({_id: id}, doc));
    return id;
}

function addTransaction(doc) {
    return Object.assign({_id: generateUUID()}, doc);
}

// Limpar dados anteriores do usuário
db.accounts.deleteMany({user_id: vasconcellosUserId});
db.categories.deleteMany({user_id: vasconcellosUserId});
db.transactions.deleteMany({user_id: vasconcellosUserId});
db.budgets.deleteMany({user_id: vasconcellosUserId});
db.goals.deleteMany({user_id: vasconcellosUserId});

// Criar categorias
const categories = {};
const categoryNames = [
    {key: "receita_salario", name: "Salário", type: "income"},
    {key: "receita_freelance", name: "Freelance", type: "income"},
    {key: "receita_investimento", name: "Renda de Investimentos", type: "income"},
    {key: "receita_bonus", name: "Bônus e Comissões", type: "income"},
    {key: "alimentacao", name: "Alimentação", type: "expense"},
    {key: "transporte", name: "Transporte", type: "expense"},
    {key: "moradia", name: "Moradia", type: "expense"},
    {key: "servicos", name: "Serviços", type: "expense"},
    {key: "saude", name: "Saúde", type: "expense"},
    {key: "educacao", name: "Educação", type: "expense"},
    {key: "tech", name: "Tecnologia", type: "expense"},
    {key: "lazer", name: "Lazer", type: "expense"},
    {key: "investimentos", name: "Investimentos", type: "expense"}
];

categoryNames.forEach(cat => {
    categories[cat.key] = insertCategory({
        user_id: vasconcellosUserId,
        name: cat.name,
        type: cat.type,
        description: cat.name,
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    });
});

// Criar contas
const accounts = {
    nubank: insertAccount({
        user_id: vasconcellosUserId,
        name: "Conta Corrente Nubank",
        type: "checking",
        currency: "BRL",
        balance: 28500.75,
        description: "Conta corrente principal - Nubank",
        created_at: new Date("2022-01-15"),
        updated_at: new Date("2024-10-25")
    }),
    inter: insertAccount({
        user_id: vasconcellosUserId,
        name: "Conta Poupança Inter",
        type: "savings",
        currency: "BRL",
        balance: 45000.00,
        description: "Reserva de emergência",
        created_at: new Date("2022-03-01"),
        updated_at: new Date("2024-10-25")
    }),
    itau: insertAccount({
        user_id: vasconcellosUserId,
        name: "Cartão de Crédito Itaú",
        type: "credit",
        currency: "BRL",
        balance: -3200.50,
        description: "Cartão de crédito Itaú",
        created_at: new Date("2022-02-10"),
        updated_at: new Date("2024-10-25")
    })
};

// Criar transações recentes (últimos 30 dias)
const transactions = [];
const now = new Date();
const startDate = new Date(now.getTime() - (30 * 24 * 60 * 60 * 1000));

for (let i = 0; i < 50; i++) {
    const randomTime = startDate.getTime() + Math.random() * (now.getTime() - startDate.getTime());
    const occurredAt = new Date(randomTime);
    
    const isIncome = Math.random() < 0.3;
    const categoryKeys = Object.keys(categories);
    const categoryId = categories[categoryKeys[Math.floor(Math.random() * categoryKeys.length)]];
    
    const accountKeys = Object.keys(accounts);
    const accountId = accounts[accountKeys[Math.floor(Math.random() * accountKeys.length)]];
    
    let amount;
    if (isIncome) {
        amount = Math.floor(Math.random() * 19500) + 500;
    } else {
        amount = -(Math.floor(Math.random() * 1490) + 10);
    }
    
    transactions.push(addTransaction({
        user_id: vasconcellosUserId,
        account_id: accountId,
        category_id: categoryId,
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
    }));
}

db.transactions.insertMany(transactions);

// Funções para criar budgets e goals com UUID
function addBudget(doc) {
    return Object.assign({_id: generateUUID()}, doc);
}

function addGoal(doc) {
    return Object.assign({_id: generateUUID()}, doc);
}

// Criar orçamentos
db.budgets.insertMany([
    addBudget({
        user_id: vasconcellosUserId,
        category_id: categories.alimentacao,
        amount: 1000,
        currency: "BRL",
        period: "monthly",
        period_start: new Date("2024-10-01"),
        period_end: new Date("2024-10-31"),
        spent: 750,
        created_at: new Date("2024-10-01"),
        updated_at: new Date("2024-10-25"),
        alert_percent: 80
    }),
    addBudget({
        user_id: vasconcellosUserId,
        category_id: categories.transporte,
        amount: 400,
        currency: "BRL",
        period: "monthly",
        period_start: new Date("2024-10-01"),
        period_end: new Date("2024-10-31"),
        spent: 380,
        created_at: new Date("2024-10-01"),
        updated_at: new Date("2024-10-25"),
        alert_percent: 80
    }),
    addBudget({
        user_id: vasconcellosUserId,
        category_id: categories.investimentos,
        amount: 2000,
        currency: "BRL",
        period: "monthly",
        period_start: new Date("2024-10-01"),
        period_end: new Date("2024-10-31"),
        spent: 1500,
        created_at: new Date("2024-10-01"),
        updated_at: new Date("2024-10-25"),
        alert_percent: 80
    })
]);

// Criar metas
db.goals.insertMany([
    addGoal({
        user_id: vasconcellosUserId,
        name: "Reserva de Emergência",
        target_amount: 90000,
        current_amount: 45000,
        currency: "BRL",
        deadline: new Date("2025-06-30"),
        status: "active",
        description: "Reserva de emergência para 6 meses",
        created_at: new Date("2022-01-15"),
        updated_at: new Date("2024-10-25")
    }),
    addGoal({
        user_id: vasconcellosUserId,
        name: "Viagem para Europa",
        target_amount: 25000,
        current_amount: 12000,
        currency: "BRL",
        deadline: new Date("2025-07-31"),
        status: "active",
        description: "Viagem de 3 semanas pela Europa",
        created_at: new Date("2023-06-01"),
        updated_at: new Date("2024-10-25")
    }),
    addGoal({
        user_id: vasconcellosUserId,
        name: "Casa Própria",
        target_amount: 150000,
        current_amount: 45000,
        currency: "BRL",
        deadline: new Date("2026-12-31"),
        status: "active",
        description: "Entrada para imóvel próprio",
        created_at: new Date("2023-01-01"),
        updated_at: new Date("2024-10-25")
    })
]);

print("\n✅ Dados do Vasconcellos criados com sucesso!");
print("📊 RESUMO:");
print("  Categories: " + Object.keys(categories).length);
print("  Accounts: " + Object.keys(accounts).length);
print("  Transactions: " + transactions.length);
print("  Budgets: 3");
print("  Goals: 3");

// Criar dados para o usuário Teste
// testeUserId já foi definido acima no início do script

db.accounts.deleteMany({user_id: testeUserId});
db.categories.deleteMany({user_id: testeUserId});
db.transactions.deleteMany({user_id: testeUserId});
db.budgets.deleteMany({user_id: testeUserId});
db.goals.deleteMany({user_id: testeUserId});

// Categorias do Teste
const testeCategories = {
    receita_salario: insertCategory({
        user_id: testeUserId,
        name: "Salário",
        type: "income",
        description: "Salário mensal",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }),
    alimentacao: insertCategory({
        user_id: testeUserId,
        name: "Alimentação",
        type: "expense",
        description: "Restaurantes e mercado",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }),
    saude: insertCategory({
        user_id: testeUserId,
        name: "Saúde",
        type: "expense",
        description: "Planos e consultas",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }),
    investimentos: insertCategory({
        user_id: testeUserId,
        name: "Investimentos",
        type: "expense",
        description: "Aplicações financeiras",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    })
};

// Contas do Teste
const testeAccounts = {
    cc_principal: insertAccount({
        user_id: testeUserId,
        name: "Conta Corrente Inter",
        type: "checking",
        currency: "BRL",
        balance: 25000.00,
        description: "Conta principal",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2024-10-25")
    })
};

// Transações recentes do Teste
const testeTransactions = [];
for (let i = 0; i < 30; i++) {
    const randomTime = startDate.getTime() + Math.random() * (now.getTime() - startDate.getTime());
    const occurredAt = new Date(randomTime);
    
    const isIncome = Math.random() < 0.4;
    const categoryId = isIncome ? testeCategories.receita_salario : 
                      (Math.random() < 0.5 ? testeCategories.alimentacao : testeCategories.saude);
    
    let amount;
    if (isIncome) {
        amount = Math.floor(Math.random() * 15000) + 1000;
    } else {
        amount = -(Math.floor(Math.random() * 1000) + 50);
    }
    
    testeTransactions.push(addTransaction({
        user_id: testeUserId,
        account_id: testeAccounts.cc_principal,
        category_id: categoryId,
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
    }));
}

db.transactions.insertMany(testeTransactions);

print("\n✅ Dados do Teste criados com sucesso!");

print("\n🎉 DEPLOY COMPLETO!");
print("📊 RESUMO FINAL:");
print("  Users: " + db.users.countDocuments());
print("  Categories: " + db.categories.countDocuments());
print("  Accounts: " + db.accounts.countDocuments());
print("  Transactions: " + db.transactions.countDocuments());
print("  Budgets: " + db.budgets.countDocuments());
print("  Goals: " + db.goals.countDocuments());

