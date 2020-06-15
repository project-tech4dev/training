"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.login = exports.authService = exports.ASClass = void 0;
const crypto_1 = __importDefault(require("crypto"));
const bank_1 = require("../models/bank");
class ASClass {
    constructor() {
        this.session = {};
        this.allowed = (req, permission) => {
            let ret = false;
            switch (permission) {
                case "accounts_list":
                case "accounts_create":
                case "users_list":
                    ret = req.userRole === bank_1.Role.BankManager;
                    break;
                case "accounts_credit":
                case "accounts_debit":
                    ret = req.userRole !== bank_1.Role.BankManager;
                    break;
                default:
                    ret = true;
                    break;
            }
            return ret;
        };
        this.isAuth = (username, password) => {
            let repl = true;
            const u = bank_1.bankService.findUserByUserName(username);
            if (u) {
                if (u.password !== password) {
                    repl = false;
                }
            }
            else {
                repl = false;
            }
            return repl;
        };
        this.getToken = (username, password) => {
            const s = crypto_1.default
                .createHash("sha256")
                .update(password)
                .digest("hex");
            return s;
        };
        this.createSession = (token, user) => {
            this.session[token] = user;
        };
        this.sessionExists = (token) => {
            return !!this.session[token];
        };
        this.getSessionUser = (token) => {
            return this.session[token];
        };
    }
}
exports.ASClass = ASClass;
exports.authService = new ASClass();
exports.login = (req, res) => {
    const username = req.body.username;
    const password = req.body.password;
    if (exports.authService.isAuth(username, password)) {
        const token = exports.authService.getToken(username, password);
        const user = bank_1.bankService.findUserByUserName(username);
        exports.authService.createSession(token, user);
        res.status(200).json({ id: user.userid, token });
    }
    else {
        res.status(500).send("Username or Password don't match");
    }
};
//# sourceMappingURL=auth.js.map