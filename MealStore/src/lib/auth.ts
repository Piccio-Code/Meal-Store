import { betterAuth } from "better-auth";
import { Pool } from "pg";

export const auth = betterAuth({
    database: new Pool({
        connectionString: process.env.DATABASE_URL, // Inserisci la tua stringa di connessione qui o nel file .env
    }),
})