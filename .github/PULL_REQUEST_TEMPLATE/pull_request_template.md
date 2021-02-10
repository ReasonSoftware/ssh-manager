## Description

Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context. List any dependencies that are required for this change.

Fixes # (issue)

## Type of change

Please delete options that are not relevant.

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

## How Has This Been Tested

Please describe the tests that you ran to verify your changes. Provide instructions so we can reproduce. Please also list any relevant details for your test configuration

- [ ] Test A
- [ ] Test B

#### Test Central Configuration

<details><summary>Clich Here to Expand</summary>

```json
{
    "users": {
        "user.1": "ssh-rsa AAA...",
        "user.2": "ssh-rsa AAA...",
        "user.3": "ssh-rsa AAA...",
        "user.4": "ssh-rsa AAA...",
        "user.5": "ssh-rsa AAA...",
        "user.6": "ssh-rsa AAA..."
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

#### Test Server Configuration

<details><summary>Clich Here to Expand</summary>

```yaml
secret_name: ssh-manager
groups:
  - devops
```

</details>

# Checklist

- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
