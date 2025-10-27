// Script para criar 1.000 transações para o Vasconcellos
db = db.getSiblingDB('financial-control');

// Buscar dados do Vasconcellos
const user = db.users.findOne({email: "vasconcellos@gmail.com"});
const vasconcellosUserId = user._id;

// Buscar contas e categorias existentes
const accounts = db.accounts.find({user_id: vasconcellosUserId}).toArray();
const categories = db.categories.find({user_id: vasconcellosUserId}).toArray();

print("📊 Criando 1.000 transações para: " + user.email);
print("Contas disponíveis: " + accounts.length);
print("Categorias disponíveis: " + categories.length);

// Criar 1.000 transações aleatórias
const transactions = [];
const categoryIds = categories.map(c => c._id.toString());
const incomeCategories = categories.filter(c => c.type === "income").map(c => c._id.toString());
const expenseCategories = categories.filter(c => c.type === "expense").map(c => c._id.toString());

// Tipos de descrições
const descriptions = {
    income: [
        "Salário mensal",
        "Freelance - Desenvolvimento",
        "Freelance - Consultoria",
        "Dividendo de investimentos",
        "Bônus trimestral",
        "Freelance - Design",
        "Juros de aplicação",
        "Venda de produto",
        "Reembolso",
        "Comissão de vendas"
    ],
    expense: [
        "Supermercado",
        "Uber",
        "Restaurante",
        "Farmácia",
        "Gasolina",
        "Cinema",
        "Curso online",
        "Netflix",
        "Aluguel",
        "Condomínio",
        "Internet + Celular",
        "Academia",
        "Livro",
        "Medicamento",
        "Ingresso show"
    ]
};

// Função para gerar data aleatória nos últimos 3 anos
function randomDate() {
    const start = new Date("2022-01-01");
    const end = new Date("2024-10-25");
    const time = start.getTime() + Math.random() * (end.getTime() - start.getTime());
    return new Date(time);
}

// Criar 1.000 transações (30% receitas, 70% despesas)
for (let i = 0; i < 1000; i++) {
    const isIncome = Math.random() < 0.3;
    const categoryIdsList = isIncome ? incomeCategories : expenseCategories;
    
    // Escolher categoria aleatória
    const categoryIndex = Math.floor(Math.random() * categoryIdsList.length);
    const categoryId = categoryIdsList[categoryIndex];
    
    // Escolher conta aleatória
    const accountIndex = Math.floor(Math.random() * accounts.length);
    const account = accounts[accountIndex];
    
    // Determinar valor
    let amount;
    if (isIncome) {
        // Receitas entre R$ 500 e R$ 20.000
        amount = Math.floor(Math.random() * 19500) + 500;
    } else {
        // Despesas entre -R$ 10 e -R$ 1.500
        amount = -(Math.floor(Math.random() * 1490) + 10);
    }
    
    // Escolher descrição aleatória
    const descList = isIncome ? descriptions.income : descriptions.expense;
    const description = descList[Math.floor(Math.random() * descList.length)];
    
    const occurredAt = randomDate();
    
    transactions.push({
        user_id: vasconcellosUserId,
        account_id: account._id.toString(),
        category_id: categoryId,
        amount: amount,
        currency: "BRL",
        description: description + " #" + (i + 1),
        occurred_at: occurredAt,
        status: "completed",
        notes: "",
        tags: isIncome ? ["receita", "aleatório"] : ["despesa", "aleatório"],
        created_at: occurredAt,
        updated_at: occurredAt,
        external_ref: "",
        metadata: {}
    });
}

// Inserir transações
const result = db.transactions.insertMany(transactions);
print("\n✅ Inseridas " + Object.keys(result.insertedIds).length + " transações do Vasconcellos com sucesso!");
print("Total de transações do Vasconcellos: " + db.transactions.countDocuments({user_id: vasconcellosUserId}));

// Criar 1.000 transações para o usuário Teste
const testeUser = db.users.findOne({email: "teste@gmail.com"});
if (testeUser) {
    const testeUserId = testeUser._id;
    const testeAccounts = db.accounts.find({user_id: testeUserId}).toArray();
    const testeCategories = db.categories.find({user_id: testeUserId}).toArray();
    
    const testeIncomeCategories = testeCategories.filter(c => c.type === "income");
    const testeExpenseCategories = testeCategories.filter(c => c.type === "expense");
    
    print("\n📊 Criando 1.000 transações para o usuário Teste");
    
    const testeTransactions = [];
    
    for (let i = 0; i < 1000; i++) {
        const isIncome = Math.random() < 0.3;
        const categoryList = isIncome ? testeIncomeCategories : testeExpenseCategories;
        const categoryIndex = Math.floor(Math.random() * categoryList.length);
        const category = categoryList[categoryIndex];
        
        const accountIndex = Math.floor(Math.random() * testeAccounts.length);
        const account = testeAccounts[accountIndex];
        
        let amount;
        if (isIncome) {
            amount = Math.floor(Math.random() * 19500) + 500;
        } else {
            amount = -(Math.floor(Math.random() * 1490) + 10);
        }
        
        const startDate = new Date("2022-01-01");
        const endDate = new Date("2024-10-25");
        const occurredAt = new Date(startDate.getTime() + Math.random() * (endDate.getTime() - startDate.getTime()));
        
        testeTransactions.push({
            user_id: testeUserId,
            account_id: account._id.toString(),
            category_id: category._id.toString(),
            amount: amount,
            currency: "BRL",
            description: (isIncome ? "Receita" : "Despesa") + " #" + (i + 1),
            occurred_at: occurredAt,
            status: "completed",
            notes: "",
            tags: isIncome ? ["receita", "aleatório"] : ["despesa", "aleatório"],
            created_at: occurredAt,
            updated_at: occurredAt,
            external_ref: "",
            metadata: {}
        });
    }
    
    const testeResult = db.transactions.insertMany(testeTransactions);
    print("✅ Inseridas " + Object.keys(testeResult.insertedIds).length + " transações do Teste com sucesso!");
    print("Total de transações do Teste: " + db.transactions.countDocuments({user_id: testeUserId}));
} else {
    print("⚠️ Usuário Teste não encontrado. Execute seed_robust_data.js primeiro.");
}

