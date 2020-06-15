"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (Object.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const body_parser_1 = __importDefault(require("body-parser"));
const express_1 = __importDefault(require("express"));
const accounts = __importStar(require("./controllers/accounts"));
const authController = __importStar(require("./controllers/auth"));
const auth_1 = require("./controllers/auth");
const users = __importStar(require("./controllers/users"));
const app = express_1.default();
// const port = 8080 || process.env.PORT;
app.use((req, res, next) => {
    const bearer = "Bearer ";
    if (req.headers.authorization) {
        const auth = req.headers.authorization;
        if (auth.length > bearer.length) {
            const token = auth.slice(bearer.length);
            const user = auth_1.authService.getSessionUser(token);
            if (user) {
                req.userID = user.userid;
                req.userRole = user.role;
            }
        }
    }
    next();
});
app.set("port", process.env.PORT || 9765);
app.use(body_parser_1.default.json());
app.use(body_parser_1.default.urlencoded({ extended: true }));
app.get("/accounts", accounts.list);
app.post("/accounts", accounts.create);
app.get("/accounts/:id", accounts.get);
app.post("/accounts/credit", accounts.credit);
app.post("/accounts/debit", accounts.debit);
app.get("/users", users.list);
app.post("/users", users.create);
app.get("/users/:id", users.get);
app.post("/login", authController.login);
app.listen(app.get("port"), () => {
    // tslint:disable-next-line:no-console
    console.log(`server started at http://localhost:${app.get("port")}`);
});
//# sourceMappingURL=index.js.map