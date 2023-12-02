# Changelog

## [v0.8.0](https://github.com/k1LoW/rl/compare/v0.7.1...v0.8.0) - 2023-11-28
### New Features 🎉
- The domain that has hit the rate limit is what I want to know. by @pyama86 in https://github.com/k1LoW/rl/pull/39

## [v0.7.1](https://github.com/k1LoW/rl/compare/v0.7.0...v0.7.1) - 2023-11-24
### Fix bug 🐛
- should check error before access rule by @pyama86 in https://github.com/k1LoW/rl/pull/37

## [v0.7.0](https://github.com/k1LoW/rl/compare/v0.6.2...v0.7.0) - 2023-10-31
### Breaking Changes 🛠
- Review Limiter interface by @k1LoW in https://github.com/k1LoW/rl/pull/35

## [v0.6.2](https://github.com/k1LoW/rl/compare/v0.6.1...v0.6.2) - 2023-10-30
### Other Changes
- Add benchmark using cachegrind by @k1LoW in https://github.com/k1LoW/rl/pull/28
- Update action by @k1LoW in https://github.com/k1LoW/rl/pull/30
- Show `cg_diff` on GitHub Actions by @k1LoW in https://github.com/k1LoW/rl/pull/31
- Change scope with newLimitError by @pyama86 in https://github.com/k1LoW/rl/pull/32

## [v0.6.1](https://github.com/k1LoW/rl/compare/v0.6.0...v0.6.1) - 2023-09-15
### Other Changes
- In the case of multiple limiters, cancel unnecessary processing if one limiter exceeds its limit. by @k1LoW in https://github.com/k1LoW/rl/pull/21
- Show benchmark in pull request using octocov by @k1LoW in https://github.com/k1LoW/rl/pull/23
- Freeze benchtime by @k1LoW in https://github.com/k1LoW/rl/pull/24
- Close request body when response error in rl by @k1LoW in https://github.com/k1LoW/rl/pull/25
- Revert "Close request body when response error in rl" by @k1LoW in https://github.com/k1LoW/rl/pull/26
- Add gostyle-action by @k1LoW in https://github.com/k1LoW/rl/pull/27

## [v0.6.0](https://github.com/k1LoW/rl/compare/v0.5.2...v0.6.0) - 2023-08-28
### Breaking Changes 🛠
- Provide Limiter with a feature to ignore the next and following Limiters. by @k1LoW in https://github.com/k1LoW/rl/pull/20
### Fix bug 🐛
- Should put a default value in the statuscode. by @pyama86 in https://github.com/k1LoW/rl/pull/15
- WriteHeader() is called before Write(). Also, it is not possible to write the header twice. by @k1LoW in https://github.com/k1LoW/rl/pull/18
### Other Changes
- Should not set `X-RateLimit-*` headers when no limit. by @k1LoW in https://github.com/k1LoW/rl/pull/19

## [v0.5.2](https://github.com/k1LoW/rl/compare/v0.5.1...v0.5.2) - 2023-08-28
### Other Changes
- No non-essential allocations by @pyama86 in https://github.com/k1LoW/rl/pull/14

## [v0.5.1](https://github.com/k1LoW/rl/compare/v0.5.0...v0.5.1) - 2023-08-27
### Breaking Changes 🛠
- Fix LimitError handling by @k1LoW in https://github.com/k1LoW/rl/pull/11

## [v0.5.0](https://github.com/k1LoW/rl/compare/v0.4.0...v0.5.0) - 2023-08-27
### Breaking Changes 🛠
- More information received when rate limits are exceeded. by @k1LoW in https://github.com/k1LoW/rl/pull/9

## [v0.4.0](https://github.com/k1LoW/rl/compare/v0.3.0...v0.4.0) - 2023-08-27
### Breaking Changes 🛠
- Remove Limiter.Name by @k1LoW in https://github.com/k1LoW/rl/pull/7
### New Features 🎉
- If reqLimit is negative, it means no limit. by @k1LoW in https://github.com/k1LoW/rl/pull/8

## [v0.3.0](https://github.com/k1LoW/rl/compare/v0.2.0...v0.3.0) - 2023-08-27
### Breaking Changes 🛠
- Re-change the interface of Limiter. by @k1LoW in https://github.com/k1LoW/rl/pull/5

## [v0.2.0](https://github.com/k1LoW/rl/compare/v0.1.0...v0.2.0) - 2023-08-27
### Breaking Changes 🛠
- Change the interface of Limiter by @k1LoW in https://github.com/k1LoW/rl/pull/3

## [v0.1.0](https://github.com/k1LoW/rl/commits/v0.1.0) - 2023-08-27
