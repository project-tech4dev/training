import bodyParser from "body-parser";
import express from "express";
import * as accounts from "./controllers/accounts";
import * as authController from "./controllers/auth";
import { authService } from "./controllers/auth";
import * as users from "./controllers/users";

export interface Request {
  userRole?: string;
  userID?: string;
}

const app = express();
// const port = 8080 || process.env.PORT;

app.use((req, res, next) => {
  const bearer = "Bearer ";
  if (req.headers.authorization) {
    const auth = req.headers.authorization;
    if (auth.length > bearer.length) {
      const token = auth.slice(bearer.length);
      const user = authService.getSessionUser(token);

      if (user) {
        req.userID = user.userid;
        req.userRole = user.role;
      }
    }
  }
  next();
});
app.set("port", process.env.PORT || 9765);
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

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
