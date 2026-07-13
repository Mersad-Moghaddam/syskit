# Security Policy

> How to report vulnerabilities in SysKit, and what to expect when you do.

---

## Overview

SysKit is a read-only Linux inspection tool. It does not modify system state, write to `/proc` or `/sys`, or perform any privileged mutation of the host. It runs with the privileges of the invoking user and reads native kernel interfaces to present system data.

Because SysKit reads and displays system information, its primary security concern is not the modification of a system but the **exposure of sensitive data** — process arguments, environment-derived details, network state, and other host information that may appear in its output. We take this responsibility seriously and welcome reports that help us keep SysKit safe to run and safe to share output from.

---

## Supported Versions

SysKit is pre-1.0 and under active development. During this phase, **only the latest released version** receives security updates. There is no long-term support for older pre-release versions — users are expected to upgrade to the most recent release.

| Version | Supported          |
|---------|--------------------|
| Latest release (`0.x`) | :white_check_mark: |
| Older `0.x` releases   | :x:                |
| Unreleased `main`      | :x: (best effort)  |

Once SysKit reaches `1.0`, this policy will be revised to define supported release lines under Semantic Versioning.

---

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues, pull requests, or discussions.**

Instead, report them privately by email to:

**security@syskit.dev**

To help us triage and resolve the issue quickly, please include as much of the following as you can:

- A clear description of the vulnerability and its impact
- The SysKit version (`syskit --version`) and how it was installed
- The Linux distribution and kernel version (`uname -a`)
- Exact steps to reproduce, including the command invoked (e.g. `syskit process --json`)
- Any relevant output, logs, or proof-of-concept — with sensitive data redacted
- Your assessment of severity and, if known, any suggested remediation

If you wish to encrypt your report, mention this in an initial email and we will coordinate a secure channel.

---

## Response Timeline

We aim to handle every report on the following schedule. These are targets, not guarantees, and we will communicate if a report requires more time.

| Stage | Target |
|---|---|
| **Acknowledgement** | Within 3 business days of receipt |
| **Initial triage & severity assessment** | Within 7 business days |
| **Fix or mitigation plan** | Within 30 days for confirmed issues |
| **Coordinated public disclosure** | After a fix is available, by agreement with the reporter |

Throughout the process we will keep you informed of progress and let you know when the issue is resolved.

---

## Disclosure Policy

SysKit follows a **coordinated disclosure** model:

1. You report the vulnerability privately to security@syskit.dev.
2. We confirm the issue, assess its impact, and develop a fix.
3. We prepare a release and, where appropriate, a security advisory.
4. We publicly disclose the vulnerability — crediting the reporter unless anonymity is requested — once a fixed release is available and users have a reasonable window to upgrade.

We ask that reporters give us a reasonable opportunity to address an issue before any public disclosure. We will not pursue or support legal action against researchers who report vulnerabilities in good faith and in accordance with this policy.

---

## Scope

SysKit is a read-only, unprivileged, Linux-only inspection tool. The following guidance clarifies what we consider in and out of scope.

### Plugin Trust Boundary

Plugins are user-installed executable code and inherit the invoking user's
permissions. Discovery and inspection never execute plugins. `plugins run`
must be explicitly invoked, refuses world-writable discovery directories and
out-of-directory executable paths, checks API compatibility, and enforces a
timeout. These checks reduce accidental exposure; they are not a sandbox.

### In Scope

- Vulnerabilities in SysKit's own code (the `syskit` binary and its libraries)
- Improper handling of untrusted data read from `/proc`, `/sys`, Netlink, or cgroups that leads to a crash, memory-safety issue, or incorrect output
- Unintended exposure of sensitive information in SysKit's output beyond what the invoking user could already access
- Injection or escaping flaws in output formatters (table, JSON, YAML) that could mislead consumers of that output
- Supply-chain issues in SysKit's declared dependencies that materially affect the binary

### Out of Scope

- Information that SysKit displays which the invoking user is already privileged to read — SysKit does not, and is not intended to, add an access-control layer over the kernel. Reading sensitive data as a user who already has permission to read it is expected behavior.
- Vulnerabilities in the Linux kernel, the C library, container runtimes, or other software SysKit inspects
- Attacks requiring root or an already-compromised host to be exploitable
- Denial of service achieved only by supplying SysKit with a deliberately hostile local environment the attacker already controls
- Issues in unsupported (older pre-1.0) versions

> **A note on output:** Because SysKit surfaces real system data, its output may contain sensitive information (such as process command lines or network details). Redirecting, logging, or sharing SysKit output can expose that data to others. This is an inherent property of an inspection tool, not a defect — but we welcome reports where SysKit exposes information the invoking user should *not* have been able to read.

---

*This security policy is a living document and will be updated as SysKit matures toward its 1.0 release. When in doubt about whether an issue qualifies, email security@syskit.dev and we will help.*
