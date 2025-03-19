# fansischmancy

Checks whether your CLI is outputting accessible color codes (i.e. 4-bit ANSI).

No fansischmancy 256 or truecolours thank you.

## Build and Run

```sh
make
<cmd> | ./fansi
```

Make sure that the command you run outputs colors when piped.

## `gh issue list --repo cli/cli`

Unable to tell that colors are not using 4-bit:

<img width="2151" alt="image" src="https://github.com/user-attachments/assets/22f2e0bd-bdf6-4cfe-8957-9018b541e48a" />

## `GH_FORCE_TTY=true gh issue list --repo cli/cli | fansi`

Very clearly you dun goofed:

<img width="2129" alt="image" src="https://github.com/user-attachments/assets/512ec7cb-ec02-4ff2-ac56-b8b004eb0f55" />

## How does it work?

<img width="390" alt="image" src="https://github.com/user-attachments/assets/23c5a27b-4930-45c1-971f-96a5c113e597" />

AI wrote it. After writing:

```go
type NonAnsiDetectionWriter struct {}
```

I've only caught glances of the code and when I saw there were regexes I shut my eyes. You certainly shouldn't trust this project with anything.
