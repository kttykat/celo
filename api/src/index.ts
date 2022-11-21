import "dotenv/config";
import express from "express";
import { O, Repo, RepoResponse } from "./types";
import chalk from "chalk";
import axios from "axios";
const app = express();
const { PORT, GITHUB_TOKEN } = process.env as unknown as O<string>;

const log = (msg: string) => {
  console.log(`${chalk.magentaBright(`celo api`)} ${chalk.gray("-")} ${msg}`);
};

app.get(`/`, (req, res) => {
  res.redirect(`https://github.com/kttykat/celo`);
});

app.get(`/:user`, async (req, res) => {
  const { user } = req.params;
  const api_url = `https://api.github.com/users/${user}/repos`;
  try {
    const h = (await axios
      .get(api_url, {
        headers: {
          authorization: `${GITHUB_TOKEN}`,
        },
      })
      .then((d) => d.data)) as Repo[];
    log(`Request made for user "${user}" found ${h.length} repos`);
    return res.send(
      h.map((z) => {
        return {
          description: z.description,
          full_name: z.full_name,
          html_url: z.html_url,
          name: z.name,
        };
      }) as RepoResponse[]
    );
  } catch (e) {
    log(`Failed to make request for user "${user}"`)
    return res.send({
      ok: false,
      data: `Failed to locate data for user: ${user}`,
    });
  }
});

app.listen(PORT, () => {
  log(`Up and running port: ${PORT}`);
});
