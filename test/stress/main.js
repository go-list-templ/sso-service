import {check} from 'k6'
import {healthz} from "./http/healthz.js"

const diagnosticUrl = 'http://app:8081'

export const options = {
    scenarios: {
        healthz_http: {
            executor: 'constant-arrival-rate',
            exec: 'runHealthz',
            rate: 2,
            timeUnit: '5s',
            duration: '1m',
            preAllocatedVUs: 1,
            maxVUs: 1,
        },
    },
    thresholds: {
        'http_req_duration{scenario:healthz_http}': ['p(95) < 500'],
        'checks': ['rate >= 0.9']
    },
    summaryTrendStats: ['min', 'max', 'p(95)', 'p(99)', 'count'],
}

export function runHealthz() {
    const response = healthz(diagnosticUrl)

    check(response, {
        'healthz is 200': (r) => r.status === 200,
    })
}