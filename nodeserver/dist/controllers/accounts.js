"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.debit = exports.credit = exports.get = exports.create = exports.list = void 0;
const bank_1 = require("../models/bank");
const auth_1 = require("./auth");
exports.list = (req, res) => {
    if (!auth_1.authService.allowed(req, "accounts_list")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const userlist = bank_1.bankService.accountList();
    res.status(200).json({ accounts: userlist });
};
exports.create = (req, res) => {
    if (!auth_1.authService.allowed(req, "accounts_create")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const userid = req.body.userid;
    const balance = req.body.balance;
    const u = bank_1.bankService.findUserByID(userid);
    if (u === null) {
        res.status(500).json({ errorcode: 500, errormessage: "No such user" });
        return;
    }
    const account = bank_1.bankService.createAccount(userid, balance);
    res.status(200).json({ accountid: account.id, balance: account.balance });
};
exports.get = (req, res) => {
    const acctid = req.params.id;
    if (!auth_1.authService.allowed(req, "accounts_get")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const acct = bank_1.bankService.findAccountByID(acctid);
    if (acct === null) {
        res.status(500).json({ errorcode: 500, errormessage: "No such account" });
        return;
    }
    if (req.userRole === bank_1.Role.User && req.userID !== acct.userid) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    res.status(200).json({
        accountid: acct.id,
        balance: acct.balance,
        activity: acct.activity,
    });
};
exports.credit = (req, res) => {
    const accountid = req.body.accountid;
    const amount = req.body.amount;
    if (!auth_1.authService.allowed(req, "accounts_credit")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const a = bank_1.bankService.findAccountByID(accountid);
    if (a === null) {
        res.status(500).json({ errorcode: 500, errormessage: "No such account" });
        return;
    }
    if (req.userRole === bank_1.Role.User && req.userID !== a.userid) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const balance = bank_1.bankService.credit(a, amount);
    res.status(200).json({ accountid: a.id, balance });
};
exports.debit = (req, res) => {
    const accountid = req.body.accountid;
    const amount = req.body.amount;
    if (!auth_1.authService.allowed(req, "accounts_debit")) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    const a = bank_1.bankService.findAccountByID(accountid);
    if (a === null) {
        res.status(500).json({ errorcode: 500, errormessage: "No such account" });
        return;
    }
    if (req.userRole === bank_1.Role.User && req.userID !== a.userid) {
        res.status(403).json({ errorcode: 403, errormessage: "Permission Denied" });
        return;
    }
    if (a.balance < amount) {
        res.status(500).json({ errorcode: 500, errormessage: "Not enough funds" });
        return;
    }
    const balance = bank_1.bankService.debit(a, amount);
    res.status(200).json({ accountid: a.id, balance });
};
//# sourceMappingURL=accounts.js.map