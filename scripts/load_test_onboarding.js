import http from 'k6/http';
import { check, sleep } from 'k6';

// k6 options to simulate a 'Ramp-Up' to 50 concurrent administrators
export const options = {
    stages: [
        { duration: '30s', target: 50 }, // Ramp up to 50 virtual users
        { duration: '1m', target: 50 },  // Stay at 50 users
        { duration: '20s', target: 0 },  // Ramp down to 0
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% of onboardings must be under 200ms
    },
};

export default function () {
    // We target the provisioning-service on port 8081
    const url = 'http://localhost:8081/v1/provision';

    const res = http.post(url);

    check(res, {
        'status is 201': (r) => r.status === 201,
        'has institution_id': (r) => r.json().hasOwnProperty('institution_id'),
    });

    sleep(1); // Wait 1 second before the next clinic onboarding
}