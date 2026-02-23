import http from 'k6/http';
import { check, sleep } from 'k6';

// Read configuration from environment variables passed via CLI
const PORT = __ENV.REST_PORT || '8090';
const BASE_URL = `http://localhost:${PORT}`;

export const options = {
    stages: [
        { duration: '30s', target: 50 }, // Ramp-up to 50 virtual users
        { duration: '1m', target: 50 },  // Steady state (maintain load)
        { duration: '20s', target: 0 },  // Ramp-down to 0
    ],
    thresholds: {
        // Performance requirements
        http_req_duration: ['p(95)<200'], // 95% of requests must complete under 200ms
        http_req_failed: ['rate<0.01'],   // Error rate must be below 1%
    },
};

export default function () {
    const url = `${BASE_URL}/v1/provision`;

    // Headers for the POST request
    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    // Sending an empty JSON body as the service generates data internally
    const res = http.post(url, JSON.stringify({}), params);

    // Validate response status and content
    check(res, {
        'is status 201': (r) => r.status === 201,
        'has institution_id': (r) => {
            try {
                return r.json().hasOwnProperty('institution_id');
            } catch (e) {
                return false;
            }
        },
    });

    // Pacing: 1-second delay between iterations per user
    sleep(1);
}