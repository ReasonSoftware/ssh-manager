---
name: Report a Bug
about: Create a report to help us improve
title: ''
labels: ''
assignees: anton-yurchenko

---

### Description

A clear and concise description of what the bug is

### Version

`Provide an Application Version`  

*First log message contains the version number...*

### Log

```Attach an execution log```  

*In case your security policy does not allow you to provide usernames, you may replace them with something like `user-1`/`user-2`.*

#### Sanitized Central Configuration

:warning: **Obfuscate real information** :warning:

<details><summary>Clich Here to Expand</summary>

```json
{
    "users": {
        "user.1": "AAA",
        "user.2": "BBB",
        "user.3": "CCC",
        "user.4": "DDD",
        "user.5": "EEE",
        "user.6": "FFF"
    },
    "server_groups": {
        "backend": {
            "sudoers": [
                "user.2"
            ],
            "users": [
                "user.1",
                "user.4",
                "user.5"
            ]
        },
        "poc": {
            "sudoers": [
                "user.1",
                "user.2",
                "user.4"
            ],
            "users": [
                "user.6"
            ]
        },
        "devops": {
            "sudoers": [
                "user.2"
            ],
            "users": [
                "user.3",
                "user.5"
            ]
        }
    }
}
```

</details>

#### Sanitized Server Configuration

:warning: **Obfuscate real information** :warning:

<details><summary>Clich Here to Expand</summary>

```yaml
secret_name: XXX
groups:
  - A
  - B
  - C
```

</details>

### Screenshots

If applicable, add screenshots to help explain your problem.
