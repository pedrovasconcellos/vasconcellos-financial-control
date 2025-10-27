// Script para criar massa de dados robusta com 3 anos de histórico
db = db.getSiblingDB('financial-control');

// Primeiro, vamos encontrar o user_id correto
const user = db.users.findOne({email: "vasconcellos@gmail.com"});
const vasconcellosUserId = user._id;

print("📊 Criando massa de dados robusta para: " + user.email);
print("User ID: " + vasconcellosUserId);

// Limpar dados anteriores do usuário
db.accounts.deleteMany({user_id: vasconcellosUserId});
db.categories.deleteMany({user_id: vasconcellosUserId});
db.transactions.deleteMany({user_id: vasconcellosUserId});
db.budgets.deleteMany({user_id: vasconcellosUserId});
db.goals.deleteMany({user_id: vasconcellosUserId});

// ============================================================================
// CATEGORIAS
// ============================================================================

const categories = {
    // RECEITAS
    receita_salario: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Salário",
        type: "income",
        description: "Salário mensal",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    receita_freelance: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Freelance",
        type: "income",
        description: "Trabalhos freelance",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    receita_investimento: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Renda de Investimentos",
        type: "income",
        description: "Dividendos e juros",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    receita_bonus: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Bônus e Comissões",
        type: "income",
        description: "Bônus e comissões",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    // DESPESAS
    alimentacao: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Alimentação",
        type: "expense",
        description: "Supermercado e restaurantes",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    transporte: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Transporte",
        type: "expense",
        description: "Uber, ônibus, gasolina",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    moradia: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Moradia",
        type: "expense",
        description: "Aluguel e condomínio",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    servicos: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Serviços",
        type: "expense",
        description: "Internet, celular, luz, água",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    saude: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Saúde",
        type: "expense",
        description: "Plano de saúde, medicamentos",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    educacao: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Educação",
        type: "expense",
        description: "Cursos, livros, certificações",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    tech: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Tecnologia",
        type: "expense",
        description: "Software, equipamentos",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    lazer: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Lazer",
        type: "expense",
        description: "Cinema, shows, viagens",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString(),
    
    investimentos: db.categories.insertOne({
        user_id: vasconcellosUserId,
        name: "Investimentos",
        type: "expense",
        description: "Aplicações financialiras",
        created_at: new Date("2022-01-01"),
        updated_at: new Date("2022-01-01")
    }).insertedId.toString()
};

// ============================================================================
// CONTAS - 3 BANCOS DIFERENTES
// ============================================================================

const accounts = {
    nubank: db.accounts.insertOne({
        user_id: vasconcellosUserId,
        name: "Conta Corrente Nubank",
        type: "checking",
        currency: "BRL",
        balance: 28500.75,
        description: "Conta corrente principal - Nubank",
        created_at: new Date("2022-01-15"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString(),
    
    inter: db.accounts.insertOne({
        user_id: vasconcellosUserId,
        name: "Conta Poupança Inter",
        type: "savings",
        currency: "BRL",
        balance: 45000.00,
        description: "Reserva de emergência",
        created_at: new Date("2022-03-01"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString(),
    
    itau: db.accounts.insertOne({
        user_id: vasconcellosUserId,
        name: "Cartão de Crédito Itaú",
        type: "credit",
        currency: "BRL",
        balance: -3200.50,
        description: "Cartão de crédito Itaú",
        created_at: new Date("2022-02-10"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString()
};

// ============================================================================
// TRANSAÇÕES - 3 ANOS DE HISTÓRICO (2022-2024)
// ============================================================================

const transactions = [];

// Função para gerar data aleatória no mês
function randomDateInMonth(year, month) {
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    const day = Math.floor(Math.random() * daysInMonth) + 1;
    return new Date(year, month, day);
}

// Gerar transações para cada mês de 2022 a 2024
for (let year = 2022; year <= 2024; year++) {
    for (let month = 0; month < 12; month++) {
        const monthStart = randomDateInMonth(year, month);
        
        // Transações de receita
        transactions.push({
            user_id: vasconcellosUserId,
            account_id: accounts.nubank,
            category_id: categories.receita_salario,
            amount: 15000,
            currency: "BRL",
            description: "Salário mensal",
            occurred_at: new Date(year, month, 5),
            status: "completed",
            notes: "",
            tags: ["salário", "mensal"],
            created_at: new Date(year, month, 5),
            updated_at: new Date(year, month, 5),
            external_ref: "",
            metadata: {}
        });
        
        // Freelances esporádicos (30% dos meses)
        if (Math.random() < 0.3) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.nubank,
                category_id: categories.receita_freelance,
                amount: Math.floor(Math.random() * 5000) + 1000,
                currency: "BRL",
                description: "Projeto freelance",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["freelance"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        // Dividendo de investimento (50% dos meses)
        if (Math.random() < 0.5) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.inter,
                category_id: categories.receita_investimento,
                amount: Math.floor(Math.random() * 300) + 100,
                currency: "BRL",
                description: "Dividendo de investimentos",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["dividendos", "investimento"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        // Despesas mensais fixas
        transactions.push({
            user_id: vasconcellosUserId,
            account_id: accounts.nubank,
            category_id: categories.moradia,
            amount: -1200,
            currency: "BRL",
            description: "Aluguel",
            occurred_at: new Date(year, month, 25),
            status: "completed",
            notes: "",
            tags: ["aluguel", "fixo"],
            created_at: new Date(year, month, 25),
            updated_at: new Date(year, month, 25),
            external_ref: "",
            metadata: {}
        });
        
        transactions.push({
            user_id: vasconcellosUserId,
            account_id: accounts.nubank,
            category_id: categories.moradia,
            amount: -350,
            currency: "BRL",
            description: "Condomínio",
            occurred_at: new Date(year, month, 10),
            status: "completed",
            notes: "",
            tags: ["condomínio", "fixo"],
            created_at: new Date(year, month, 10),
            updated_at: new Date(year, month, 10),
            external_ref: "",
            metadata: {}
        });
        
        transactions.push({
            user_id: vasconcellosUserId,
            account_id: accounts.nubank,
            category_id: categories.servicos,
            amount: -450,
            currency: "BRL",
            description: "Internet + Celular",
            occurred_at: new Date(year, month, 15),
            status: "completed",
            notes: "",
            tags: ["serviços", "fixo"],
            created_at: new Date(year, month, 15),
            updated_at: new Date(year, month, 15),
            external_ref: "",
            metadata: {}
        });
        
        transactions.push({
            user_id: vasconcellosUserId,
            account_id: accounts.nubank,
            category_id: categories.saude,
            amount: -850,
            currency: "BRL",
            description: "Plano de saúde",
            occurred_at: new Date(year, month, 8),
            status: "completed",
            notes: "",
            tags: ["saúde", "fixo"],
            created_at: new Date(year, month, 8),
            updated_at: new Date(year, month, 8),
            external_ref: "",
            metadata: {}
        });
        
        // Despesas variáveis (alimentação, transporte)
        const alimentacaoCount = Math.floor(Math.random() * 3) + 2;
        for (let i = 0; i < alimentacaoCount; i++) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.itau,
                category_id: categories.alimentacao,
                amount: -(Math.floor(Math.random() * 200) + 50),
                currency: "BRL",
                description: "Supermercado",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["supermercado", "alimentação"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        const transporteCount = Math.floor(Math.random() * 10) + 5;
        for (let i = 0; i < transporteCount; i++) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.itau,
                category_id: categories.transporte,
                amount: -(Math.floor(Math.random() * 50) + 15),
                currency: "BRL",
                description: "Uber",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["uber", "transporte"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        // Despesas ocasionais
        if (Math.random() < 0.3) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.itau,
                category_id: categories.educacao,
                amount: -(Math.floor(Math.random() * 1000) + 300),
                currency: "BRL",
                description: "Curso online",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["educação", "curso"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        if (Math.random() < 0.2) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.itau,
                category_id: categories.tech,
                amount: -(Math.floor(Math.random() * 500) + 100),
                currency: "BRL",
                description: "Software/Assinatura",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["tecnologia", "software"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        if (Math.random() < 0.15) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.itau,
                category_id: categories.lazer,
                amount: -(Math.floor(Math.random() * 500) + 200),
                currency: "BRL",
                description: "Cinema/Restaurante",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["lazer", "entretenimento"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
        
        // Investimentos (50% dos meses)
        if (Math.random() < 0.5) {
            transactions.push({
                user_id: vasconcellosUserId,
                account_id: accounts.inter,
                category_id: categories.investimentos,
                amount: -(Math.floor(Math.random() * 3000) + 1000),
                currency: "BRL",
                description: "Aplicação em investimentos",
                occurred_at: randomDateInMonth(year, month),
                status: "completed",
                notes: "",
                tags: ["investimento", "aplicação"],
                created_at: randomDateInMonth(year, month),
                updated_at: randomDateInMonth(year, month),
                external_ref: "",
                metadata: {}
            });
        }
    }
}

// Inserir transações
const result = db.transactions.insertMany(transactions);
print("✅ Inseridas " + result.insertedIds.length + " transações");

// ============================================================================
// ORÇAMENTOS ATUAIS
// ============================================================================

db.budgets.insertMany([
    {
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
    },
    {
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
    },
    {
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
    }
]);

// ============================================================================
// METAS FINANCEIRAS
// ============================================================================

db.goals.insertMany([
    {
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
    },
    {
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
    },
    {
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
    }
]);

print("\n📊 RESUMO DA MASSA DE DADOS DO VASCONCELLOS:");
print("   ✅ " + Object.keys(categories).length + " categorias criadas");
print("   ✅ " + Object.keys(accounts).length + " contas criadas (Nubank, Inter, Itaú)");
print("   ✅ " + result.insertedIds.length + " transações criadas (3 anos de histórico)");
print("   ✅ 3 orçamentos ativos");
print("   ✅ 3 metas financialiras");
print("\n🎉 Massa de dados do Vasconcellos criada com sucesso!");

// ============================================================================
// DADOS DO TESTE
// ============================================================================

const testeUser = db.users.findOne({email: "teste@gmail.com"});
const testeUserId = testeUser._id;

print("\n📊 Criando massa de dados para: " + testeUser.email);
print("User ID: " + testeUserId);

// Limpar dados anteriores do Teste
db.accounts.deleteMany({user_id: testeUserId});
db.categories.deleteMany({user_id: testeUserId});
db.transactions.deleteMany({user_id: testeUserId});
db.budgets.deleteMany({user_id: testeUserId});
db.goals.deleteMany({user_id: testeUserId});

// Categorias do Teste
const testeCategories = {
    receita_salario: db.categories.insertOne({
        user_id: testeUserId,
        name: "Salário",
        type: "income",
        description: "Salário mensal",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }).insertedId.toString(),
    
    receita_investimento: db.categories.insertOne({
        user_id: testeUserId,
        name: "Renda de Investimentos",
        type: "income",
        description: "Dividendos e juros",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }).insertedId.toString(),
    
    alimentacao: db.categories.insertOne({
        user_id: testeUserId,
        name: "Alimentação",
        type: "expense",
        description: "Restaurantes e mercado",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }).insertedId.toString(),
    
    saude: db.categories.insertOne({
        user_id: testeUserId,
        name: "Saúde",
        type: "expense",
        description: "Planos e consultas",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }).insertedId.toString(),
    
    educacao: db.categories.insertOne({
        user_id: testeUserId,
        name: "Educação",
        type: "expense",
        description: "Cursos e livros",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }).insertedId.toString(),
    
    investimentos: db.categories.insertOne({
        user_id: testeUserId,
        name: "Investimentos",
        type: "expense",
        description: "Aplicações financialiras",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2022-02-20")
    }).insertedId.toString()
};

// Contas do Teste
const testeAccounts = {
    cc_principal: db.accounts.insertOne({
        user_id: testeUserId,
        name: "Conta Corrente Inter",
        type: "checking",
        currency: "BRL",
        balance: 25000.00,
        description: "Conta principal",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString(),
    
    investimentos: db.accounts.insertOne({
        user_id: testeUserId,
        name: "Carteira de Investimentos",
        type: "savings",
        currency: "BRL",
        balance: 120000.00,
        description: "Investimentos diversos",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString(),
    
    credito: db.accounts.insertOne({
        user_id: testeUserId,
        name: "Cartão de Crédito",
        type: "credit",
        currency: "BRL",
        balance: -1500.00,
        description: "Visa Platinum",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2024-10-25")
    }).insertedId.toString()
};

// Transações do Teste (3 anos)
const testeTransactions = [];

for (let year = 2022; year <= 2024; year++) {
    for (let month = 0; month < 12; month++) {
        const monthStart = new Date(year, month, Math.floor(Math.random() * 28) + 1);
        
        // Salário mensal
        testeTransactions.push({
            user_id: testeUserId,
            account_id: testeAccounts.cc_principal,
            category_id: testeCategories.receita_salario,
            amount: 18000,
            currency: "BRL",
            description: "Salário mensal",
            occurred_at: new Date(year, month, 5),
            status: "completed",
            notes: "",
            tags: ["salário", "mensal"],
            created_at: new Date(year, month, 5),
            updated_at: new Date(year, month, 5),
            external_ref: "",
            metadata: {}
        });
        
        // Investimentos (50% dos meses)
        if (Math.random() < 0.5) {
            testeTransactions.push({
                user_id: testeUserId,
                account_id: testeAccounts.investimentos,
                category_id: testeCategories.investimentos,
                amount: -(Math.floor(Math.random() * 5000) + 2000),
                currency: "BRL",
                description: "Aplicação em investimentos",
                occurred_at: new Date(year, month, Math.floor(Math.random() * 28) + 1),
                status: "completed",
                notes: "",
                tags: ["investimento", "aplicação"],
                created_at: new Date(year, month, Math.floor(Math.random() * 28) + 1),
                updated_at: new Date(year, month, Math.floor(Math.random() * 28) + 1),
                external_ref: "",
                metadata: {}
            });
        }
        
        // Despesas ocasionais
        if (Math.random() < 0.3) {
            testeTransactions.push({
                user_id: testeUserId,
                account_id: testeAccounts.credito,
                category_id: testeCategories.saude,
                amount: -(Math.floor(Math.random() * 500) + 300),
                currency: "BRL",
                description: "Consulta médica",
                occurred_at: new Date(year, month, Math.floor(Math.random() * 28) + 1),
                status: "completed",
                notes: "",
                tags: ["saúde"],
                created_at: new Date(year, month, Math.floor(Math.random() * 28) + 1),
                updated_at: new Date(year, month, Math.floor(Math.random() * 28) + 1),
                external_ref: "",
                metadata: {}
            });
        }
    }
}

const testeResult = db.transactions.insertMany(testeTransactions);
print("✅ Inseridas " + Object.keys(testeResult.insertedIds).length + " transações do Teste");

// Orçamentos do Teste
db.budgets.insertMany([
    {
        user_id: testeUserId,
        category_id: testeCategories.investimentos,
        amount: 8000,
        currency: "BRL",
        period: "monthly",
        period_start: new Date("2024-10-01"),
        period_end: new Date("2024-10-31"),
        spent: 5000,
        created_at: new Date("2024-10-01"),
        updated_at: new Date("2024-10-25"),
        alert_percent: 80
    },
    {
        user_id: testeUserId,
        category_id: testeCategories.saude,
        amount: 1500,
        currency: "BRL",
        period: "monthly",
        period_start: new Date("2024-10-01"),
        period_end: new Date("2024-10-31"),
        spent: 900,
        created_at: new Date("2024-10-01"),
        updated_at: new Date("2024-10-25"),
        alert_percent: 80
    },
    {
        user_id: testeUserId,
        category_id: testeCategories.educacao,
        amount: 3000,
        currency: "BRL",
        period: "monthly",
        period_start: new Date("2024-10-01"),
        period_end: new Date("2024-10-31"),
        spent: 1500,
        created_at: new Date("2024-10-01"),
        updated_at: new Date("2024-10-25"),
        alert_percent: 80
    }
]);

// Metas do Teste
db.goals.insertMany([
    {
        user_id: testeUserId,
        name: "Independência Financialira",
        target_amount: 500000,
        current_amount: 120000,
        currency: "BRL",
        deadline: new Date("2028-12-31"),
        status: "active",
        description: "Meta de patrimônio para independência financialira",
        created_at: new Date("2022-02-20"),
        updated_at: new Date("2024-10-25")
    },
    {
        user_id: testeUserId,
        name: "Casa Própria",
        target_amount: 300000,
        current_amount: 95000,
        currency: "BRL",
        deadline: new Date("2026-06-30"),
        status: "active",
        description: "Entrada para imóvel próprio",
        created_at: new Date("2023-01-01"),
        updated_at: new Date("2024-10-25")
    }
]);

print("\n📊 RESUMO DA MASSA DE DADOS DO TESTE:");
print("   ✅ " + Object.keys(testeCategories).length + " categorias criadas");
print("   ✅ " + Object.keys(testeAccounts).length + " contas criadas");
print("   ✅ " + Object.keys(testeResult.insertedIds).length + " transações criadas");
print("   ✅ 3 orçamentos ativos");
print("   ✅ 2 metas financialiras");
print("\n🎉 Massa de dados do Teste criada com sucesso!");

