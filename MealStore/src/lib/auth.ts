import { betterAuth } from "better-auth";
import { Pool } from "pg";
import { jwt } from "better-auth/plugins";

// Create PostgreSQL connection pool
const pool = new Pool({
    connectionString: process.env.DATABASE_URL,
});

export const auth = betterAuth({
    database: pool,

    // JWT Plugin Configuration
    plugins: [
        jwt(),
    ],

    // Email and Password Authentication
    emailAndPassword: {
        enabled: true
    },

    // Base URL Configuration
    baseURL: process.env.BETTER_AUTH_URL,
});

