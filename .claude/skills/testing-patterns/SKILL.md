---
name: Testing Patterns
description: This skill should be used when the user asks to "write tests", "test strategy", "coverage", "unit test", "integration test", or needs testing guidance. Provides testing methodology and patterns.
---

# Testing Patterns and Strategy

## Test Types

| Type | Scope | Speed | When to Use |
|------|-------|-------|-------------|
| Unit | Single function/class | Fast | Business logic, utilities, transformations |
| Integration | Multiple components | Medium | API endpoints, DB operations, services |
| E2E | Full application | Slow | Critical user journeys, smoke tests |

## Test Structure

### Arrange-Act-Assert (Technical)
```
// Arrange: Set up test data
user = User.new(name: "John")
cart = ShoppingCart.new(user)

// Act: Execute code under test
total = cart.calculate_total

// Assert: Verify outcomes
assert_equal 0, total
```

### Given-When-Then (BDD)
```
given_a_user_with_an_empty_cart
when_the_user_calculates_total
then_the_total_should_be_zero
```

## Test Doubles

| Type | Purpose | Example |
|------|---------|---------|
| Stub | Canned responses | `stub(fetch_user: { id: 1 })` |
| Mock | Verify interactions | `mock.expect(:send_email, args: [...])` |
| Spy | Record calls, real behavior | `spy(Logger.new)` |
| Fake | Working lightweight impl | In-memory database |

## Naming Conventions

**Technical**: `test_[method]_[scenario]_[expected]`
```
test_calculateTotal_withEmptyCart_returnsZero
```

**BDD**: `[method]_should_[behavior]_when_[condition]`
```
calculateTotal_should_returnZero_when_cartIsEmpty
```

## Critical Practices

1. **Test happy path first**, then edge cases, then errors
2. **Edge cases**: Empty inputs, max values, nulls, zero, negatives
3. **Error cases**: Invalid inputs, network failures, timeouts
4. **Isolate tests**: Each test independent with setup/teardown
5. **One assertion per concept**: Single logical verification per test
6. **Readable names**: Tests serve as documentation

## Coverage Guidelines
- 80%+ for critical code paths
- Branch coverage more thorough than line coverage
- High coverage does not guarantee bug-free code
- Prioritize meaningful tests over metrics

## Anti-Patterns
- **Testing implementation** - Test behavior, not internals
- **Excessive mocking** - Use real implementations where practical
- **Flaky tests** - Control time, randomness, async operations
- **Slow tests** - Unit tests should run in milliseconds
- **Test interdependence** - Each test creates its own data
