"use strict";
/* tslint:disable:max-classes-per-file */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.bankService = exports.userService = exports.Role = void 0;
var Role;
(function (Role) {
    Role["BankManager"] = "BankManager";
    Role["User"] = "User";
})(Role = exports.Role || (exports.Role = {}));
const fs_1 = __importDefault(require("fs"));
const Version = 1;
class BSClass {
    constructor() {
        this.bankfile = "data/bank.json";
        this.init = () => {
            if (fs_1.default.existsSync(this.bankfile)) {
                this.bank = JSON.parse(fs_1.default.readFileSync(this.bankfile, "utf8"));
            }
            if (!this.bank || this.bank.version !== Version) {
                this.bank = {
                    version: Version,
                    accountserial: 10001,
                    accountscatalog: {},
                    nextuserid: 1000000,
                    useridcatalog: {},
                    usernamecatalog: {},
                    users: [],
                };
                // console.log(this.bank);
                const userid = exports.userService.create("Stockson Bonds", "bankmanager", "headhoncho", Role.BankManager);
                this.flush();
            }
            // console.log(this.bank);
        };
        this.nextuserid = () => {
            this.bank.nextuserid++;
            this.flush();
            return this.bank.nextuserid.toString();
        };
        this.nextaccountid = () => {
            this.bank.accountserial++;
            this.flush();
            return this.bank.accountserial.toString();
        };
        this.addUser = (u) => {
            this.bank.users.push(u);
            this.bank.useridcatalog[u.userid] = this.bank.users.length - 1;
            this.bank.usernamecatalog[u.username] = this.bank.users.length - 1;
            this.flush();
        };
        this.flush = () => {
            fs_1.default.writeFileSync(this.bankfile, JSON.stringify(this.bank));
        };
        this.findUserByID = (id) => {
            if (!this.bank.useridcatalog[id]) {
                return null;
            }
            const ix = this.bank.useridcatalog[id];
            const u = this.bank.users[ix];
            return u;
        };
        this.findUserByUserName = (uname) => {
            const ix = this.bank.usernamecatalog[uname];
            if (ix === undefined) {
                return null;
            }
            return this.bank.users[ix];
        };
        this.findAccountByID = (id) => {
            if (!this.bank.accountscatalog[id]) {
                return null;
            }
            const uid = this.bank.accountscatalog[id];
            const u = this.findUserByID(uid);
            const a = u.accounts.find((e) => e.id === id);
            return a;
        };
        this.addAccount = (userid, account) => {
            const ix = this.bank.useridcatalog[userid];
            this.bank.users[ix].accounts.push(account);
            this.bank.accountscatalog[account.id] = userid;
        };
        this.accountList = () => {
            const list = [];
            this.bank.users.forEach((u) => {
                u.accounts.forEach((a) => {
                    list.push({
                        name: u.name,
                        accountid: a.id,
                        balance: a.balance,
                        createdat: a.createdat,
                    });
                });
            });
            return list;
        };
        this.createAccount = (userid, balance) => {
            const id = this.nextaccountid();
            const t = {
                amount: balance,
                balance,
                createdat: new Date(),
                operation: "OB",
            };
            const a = {
                id,
                userid,
                balance,
                createdat: new Date(),
                activity: [t],
            };
            this.bank.accountscatalog[id] = userid;
            const u = this.findUserByID(userid); // user exists, we checked in the controller
            u.accounts.push(a);
            this.flush();
            return a;
        };
        this.credit = (account, amount) => {
            const balance = account.balance + amount;
            const t = {
                amount,
                balance,
                createdat: new Date(),
                operation: "CR",
            };
            account.balance += amount;
            account.activity.push(t);
            this.flush();
            return t.balance;
        };
        this.debit = (account, amount) => {
            const balance = account.balance - amount;
            const t = {
                amount,
                balance,
                createdat: new Date(),
                operation: "DB",
            };
            account.balance -= amount;
            account.activity.push(t);
            this.flush();
            return t.balance;
        };
    }
}
class USClass {
    constructor() {
        this.create = (fullname, username, password, role = Role.User) => {
            // console.log("kkkkkk", bankService.bank);
            const userid = exports.bankService.nextuserid();
            const u = {
                userid: userid.toString(),
                name: fullname,
                username,
                password,
                role,
                createdat: new Date(),
                accounts: [],
            };
            exports.bankService.addUser(u);
            return u.userid;
        };
        this.userList = () => {
            const list = exports.bankService.bank.users
                .map((u) => {
                if (u.role !== Role.BankManager) {
                    return {
                        id: u.userid,
                        fullname: u.name,
                        username: u.username,
                        createdat: u.createdat,
                        numaccounts: u.accounts.length,
                    };
                }
            })
                .filter((l) => l != null);
            return list;
        };
    }
}
exports.userService = new USClass();
exports.bankService = new BSClass();
exports.bankService.init();
//# sourceMappingURL=bank.js.map