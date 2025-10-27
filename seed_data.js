// Script para inserir massa de dados no MongoDB
db = db.getSiblingDB('financial-control');

const vasconcellosUserId = "local-vasconcellos-user";

// Inserir categorias do Vasconcellos
const vasconcellosCategories = {
    receita_freelance: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Freelance",
        type: "income",
        description: "Trabalhos freelance",
        created_at: new Date("2024-01-15"),
        updated_at: new Date("2024-01-15")
    }).insertedId.toString(),
    
    receita_salario: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Salário",
        type: "income",
        description: "Salário mensal",
        created_at: new Date("2024-01-15"),
        updated_at: new Date("2024-01-15")
    }).insertedId.toString(),
    
    alimentacao: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Alimentação",
        type: "expense",
        description: "Supermercado e restaurantes",
        created_at: new Date("2024-01-15"),
        updated_at: new Date("2024-01-15")
    }).insertedId.toString(),
    
    transporte: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Transporte",
        type: "expense",
        description: "Uber, ônibus, gasolina",
        created_at: new Date("2024-01-15"),
        updated_at: new Date("2024-01-15")
    }).insertedId.toString(),
    
    moradia: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Moradia",
        type: "expense",
        description: "Aluguel e condomínio",
        created_at: new Date("2024-01-15"),
        updated_at: new Date("2024-01-15")
    }).insertedId.toString()
};

// Inserir contas do Vasconcellos
const vasconcellosAccounts = {
    cc: db.accounts.insertOne({
        user_id: vasconcellosUserId,
        name: "Conta Corrente Nubank",
        type: "checking",
        currency: "BRL",
        balance: 18500.75,
        description: "Conta corrente principal",
        created_at: new Date("2024-01-15"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString()
};

// Inserir transações do Vasconcellos
db.transactions.insertMany([
    {
        user_id: vasconcellosUserId,
        account_id: vasconcellosAccounts.cc,
        category_id: vasconcellosCategories.receita_salario,
        amount: 15000,
        currency: "BRL",
        description: "Salário mensal",
        occurred_at: new Date("2024-10-05"),
        status: "completed",
        notes: "",
        tags: [],
        created_at: new Date("2024-10-05"),
        updated_at: new Date("2024-10-05"),
        external_ref: "",
        metadata: {}
    },
    {
        user_id: vasconcellosUserId,
        account_id: vasconcellosAccounts.cc,
        category_id: vasconcellosCategories.alimentacao,
        amount: -850,
        currency: "BRL",
        description: "Supermercado",
        occurred_at: new Date("2024-10-10"),
        status: "completed",
        notes: "",
        tags: [],
        created_at: new Date("2024-10-10"),
        updated_at: new Date("2024-10-10"),
        external_ref: "",
        metadata: {}
    },
    {
        user_id: vasconcellosUserId,
        account_id: vasconcellosAccounts.cc,
        category_id: vasconcellosCategories.transporte,
        amount: -350,
        currency: "BRL",
        description: "Uber",
        occurred_at: new Date("2024-10-15"),
        status: "completed",
        notes: "",
        tags: [],
        created_at: new Date("2024-10-15"),
        updated_at: new Date("2024-10-15"),
        external_ref: "",
        metadata: {}
    }
]);

print("✅ Dados do Vasconcellos inseridos com sucesso!");

