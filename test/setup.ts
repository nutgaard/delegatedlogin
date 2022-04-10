const useColor = Deno.env.get("CI") !== 'true';
const RED = useColor ? '\x1b[31m' : '';
const GREEN = useColor ? '\x1b[32m' : '';
const CYAN = useColor ? '\x1b[36m' : '';
const RESET = useColor ? '\x1b[0m' : '';

export type Retry = { retry: number; interval: number; }

function sleep(seconds: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, seconds * 1000));
}

export async function setup(name: string, exec: () => Promise<void>) {
    console.log(`${CYAN}[ UP ]${RESET} ${name}`);
    await exec();
    console.log();
}

export function retry({ retry, interval }: Retry, exec: (attempt: number) => any): () => Promise<void> {
    let count = 0;
    return async () => {
        do {
            try {
                await exec(count + 1);
                return;
            } catch (e) {
                console.log(`${RED} [KO] ${RESET} ${e}`);
            }
            count++;
            await sleep(interval);
        } while (count < retry);
        console.log();
        console.log(`${RED} [KO] ${RESET} Setup failed. Exiting...`);
        console.log();
        Deno.exit(1);
    };
}