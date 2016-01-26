/// <reference path="../typings/tsd.d.ts" />
let parser = require("tap-parser");

interface Result {
    ok: boolean;
    count: number;
    pass: number;
    fail: number;
    plan: Plan;
    failures: Failure[];
}
interface Plan {
    start: number;
    end: number;
}
interface Failure {
    ok: boolean;
    id: number;
    name: string;
    diag: Diagnostic;
}
interface Diagnostic {
    operator: string;
    expected: any;
    actual: any;
    at: string;
}

interface OutputLine {
    message: string;
    line: string;
    expected: any;
    actual: any;
}

let p = parser((result: Result) => {
    if (result.ok) {
        console.dir({});
        return;
    }

    let output: OutputLine[] = result.failures.map((failure: Failure) => {
        let message = failure.name;
        if (message.length === 0) {
            message = "<no message provided>";
        }

        return <OutputLine>{
            message: message,
            line: failure.diag.at,
            expected: failure.diag.expected.toString(),
            actual: failure.diag.actual.toString()
        };
    });
    console.log(JSON.stringify({ output: output }));
    process.exit(1);
});
process.stdin.pipe(p);
