<!--
This file is for important supplementary information that needs to be included in the next release.
This is not meant to duplicate information from PR descriptions or commit messages.
Examples might include:
- Breaking changes that need extra explanation
- Migration guides
- Complex upgrade instructions
- Important dependency changes that need special attention
-->

## Additional Release Information

This release adds a new dedicated sanitizer for BCP. This sanitizer can be enabled by setting the feature toggle flag `ENABLE_BCP_SANITIZER` to true as an environmental variable. This sanitizer is built to ensure guaranteed consistency in generating the sanitized block labels.
