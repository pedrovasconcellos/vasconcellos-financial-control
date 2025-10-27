// Script de migração para converter ObjectId em UUID (string) em todas as coleções
db = db.getSiblingDB('financial-control');

function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

function toStringId(value) {
  if (value === null || value === undefined) {
    return value;
  }
  if (typeof value === 'string') {
    return value;
  }
  if (value instanceof ObjectId) {
    return value.toString();
  }
  return String(value);
}

function convertUsers() {
  print('🔄 Convertendo usuários para UUID...');
  db.users.find().forEach((user) => {
    const oldId = user._id;
    const oldIdStr = toStringId(oldId);
    const userIdIsString = typeof oldId === 'string';

    let newId = oldIdStr;
    if (!userIdIsString) {
      newId = generateUUID();
      const newUser = Object.assign({}, user);
      newUser._id = newId;

      db.users.deleteOne({ _id: oldId });
      db.users.insertOne(newUser);
    }

    const userIdStringForRefs = toStringId(newId);
    const possibleValues = [oldId, oldIdStr];

    db.accounts.updateMany(
      { user_id: { $in: possibleValues } },
      { $set: { user_id: userIdStringForRefs } }
    );
    db.categories.updateMany(
      { user_id: { $in: possibleValues } },
      { $set: { user_id: userIdStringForRefs } }
    );
    db.transactions.updateMany(
      { user_id: { $in: possibleValues } },
      { $set: { user_id: userIdStringForRefs } }
    );
    db.budgets.updateMany(
      { user_id: { $in: possibleValues } },
      { $set: { user_id: userIdStringForRefs } }
    );
    db.goals.updateMany(
      { user_id: { $in: possibleValues } },
      { $set: { user_id: userIdStringForRefs } }
    );
  });
}

function convertAccounts() {
  print('🔄 Convertendo contas para UUID...');
  db.accounts.find().forEach((account) => {
    const oldId = account._id;
    const oldIdStr = toStringId(oldId);
    const needsConversion = typeof oldId !== 'string';
    const userIdStr = toStringId(account.user_id);

    if (!needsConversion) {
      if (account.user_id !== userIdStr) {
        db.accounts.updateOne({ _id: oldId }, { $set: { user_id: userIdStr } });
      }
      return;
    }

    const newId = generateUUID();
    const newAccount = Object.assign({}, account);
    newAccount._id = newId;
    newAccount.user_id = userIdStr;

    db.accounts.deleteOne({ _id: oldId });
    db.accounts.insertOne(newAccount);

    db.transactions.updateMany(
      { account_id: { $in: [oldId, oldIdStr] } },
      { $set: { account_id: newId } }
    );
  });
}

function convertCategories() {
  print('🔄 Convertendo categorias para UUID...');
  db.categories.find().forEach((category) => {
    const oldId = category._id;
    const oldIdStr = toStringId(oldId);
    const needsConversion = typeof oldId !== 'string';
    const userIdStr = toStringId(category.user_id);
    const parentIdStr = toStringId(category.parent_id);

    if (!needsConversion) {
      const updates = {};
      if (category.user_id !== userIdStr) {
        updates.user_id = userIdStr;
      }
      if (
        category.parent_id !== undefined &&
        category.parent_id !== null &&
        category.parent_id !== parentIdStr
      ) {
        updates.parent_id = parentIdStr;
      }
      if (Object.keys(updates).length > 0) {
        db.categories.updateOne({ _id: oldId }, { $set: updates });
      }
      return;
    }

    const newId = generateUUID();
    const newCategory = Object.assign({}, category);
    newCategory._id = newId;
    newCategory.user_id = userIdStr;
    newCategory.parent_id = parentIdStr;

    db.categories.deleteOne({ _id: oldId });
    db.categories.insertOne(newCategory);

    db.transactions.updateMany(
      { category_id: { $in: [oldId, oldIdStr] } },
      { $set: { category_id: newId } }
    );
    db.budgets.updateMany(
      { category_id: { $in: [oldId, oldIdStr] } },
      { $set: { category_id: newId } }
    );
    db.categories.updateMany(
      { parent_id: { $in: [oldId, oldIdStr] } },
      { $set: { parent_id: newId } }
    );
  });
}

function convertBudgets() {
  print('🔄 Convertendo orçamentos para UUID...');
  db.budgets.find().forEach((budget) => {
    const oldId = budget._id;
    const needsConversion = typeof oldId !== 'string';
    const userIdStr = toStringId(budget.user_id);
    const categoryIdStr = toStringId(budget.category_id);

    if (!needsConversion) {
      const updates = {};
      if (budget.user_id !== userIdStr) {
        updates.user_id = userIdStr;
      }
      if (budget.category_id !== categoryIdStr) {
        updates.category_id = categoryIdStr;
      }
      if (Object.keys(updates).length > 0) {
        db.budgets.updateOne({ _id: oldId }, { $set: updates });
      }
      return;
    }

    const newId = generateUUID();
    const newBudget = Object.assign({}, budget);
    newBudget._id = newId;
    newBudget.user_id = userIdStr;
    newBudget.category_id = categoryIdStr;

    db.budgets.deleteOne({ _id: oldId });
    db.budgets.insertOne(newBudget);
  });
}

function convertGoals() {
  print('🔄 Convertendo metas para UUID...');
  db.goals.find().forEach((goal) => {
    const oldId = goal._id;
    const needsConversion = typeof oldId !== 'string';
    const userIdStr = toStringId(goal.user_id);

    if (!needsConversion) {
      if (goal.user_id !== userIdStr) {
        db.goals.updateOne({ _id: oldId }, { $set: { user_id: userIdStr } });
      }
      return;
    }

    const newId = generateUUID();
    const newGoal = Object.assign({}, goal);
    newGoal._id = newId;
    newGoal.user_id = userIdStr;

    db.goals.deleteOne({ _id: oldId });
    db.goals.insertOne(newGoal);
  });
}

function convertTransactions() {
  print('🔄 Convertendo transações para UUID...');
  db.transactions.find().forEach((transaction) => {
    const oldId = transaction._id;
    const needsConversion = typeof oldId !== 'string';

    const userIdStr = toStringId(transaction.user_id);
    const accountIdStr = toStringId(transaction.account_id);
    const categoryIdStr = toStringId(transaction.category_id);

    if (!needsConversion) {
      const updates = {};
      if (transaction.user_id !== userIdStr) {
        updates.user_id = userIdStr;
      }
      if (transaction.account_id !== accountIdStr) {
        updates.account_id = accountIdStr;
      }
      if (transaction.category_id !== categoryIdStr) {
        updates.category_id = categoryIdStr;
      }
      if (Object.keys(updates).length > 0) {
        db.transactions.updateOne({ _id: oldId }, { $set: updates });
      }
      return;
    }

    const newId = generateUUID();
    const newTransaction = Object.assign({}, transaction);
    newTransaction._id = newId;
    newTransaction.user_id = userIdStr;
    newTransaction.account_id = accountIdStr;
    newTransaction.category_id = categoryIdStr;

    db.transactions.deleteOne({ _id: oldId });
    db.transactions.insertOne(newTransaction);
  });
}

convertUsers();
convertAccounts();
convertCategories();
convertBudgets();
convertGoals();
convertTransactions();

print('✅ Conversão concluída. Todos os IDs agora são UUIDs string.');
