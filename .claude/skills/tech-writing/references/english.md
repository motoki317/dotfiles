# English

## House style

- Use the active voice when you know the actor: "The CLI runs the worker", not "The worker is run".
- Write no contractions and no semicolons.
- Use "must" for a requirement and "can" for a capability. Keep "may", "might", and "could" for uncertainty, and "would" for a counterfactual. Write a recommendation as the fact that supports it: "X avoids Y", not "You should use X".

## Self-check

Run this command on the file before you post or save it. A hit in quoted text is not a violation.

```sh
perl -CSD -Mutf8 -ne 'if(/^\s*```/){$c=!$c;next} next if $c; s/`[^`]*`//g; print "$.: $&\n" while /\b\w+(?:n[\x27’]t|[\x27’](?:ll|re|ve|d))\b|\b(?:it|that|there|let)[\x27’]s\b|\b(?:should|would|may|might|could)\b|;/gi' FILE
```
