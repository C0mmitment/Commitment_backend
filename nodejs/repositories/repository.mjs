import util from 'node:util';
import { exec } from 'node:child_process';
const execPromise = util.promisify(exec);

const apiTestNetwork = async (ip) => {
    try {
        await execPromise(`ping -c 1 -W 1 ${ip}`);
        return true;
    } catch (err) {
        console.error('Network test failed:', err);
        return false;
    }
}

export default {
    apiTestNetwork
}