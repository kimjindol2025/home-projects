/**
 * Test Suite: Design Intent Extractor (Phase C)
 *
 * 7개 테스트 시나리오:
 * 1. 문서에서 설계 의도 추출
 * 2. 이름 규칙에서 의도 파악
 * 3. 구조에서 패턴 의도 감지
 * 4. 구현과 의도의 괴리 감지
 * 5. 의도 진화 추적
 * 6. 설계 결정 이유 분석
 * 7. 통합 설계 의도 분석
 */

const { DesignIntentExtractor } = require('./design-intent-extractor.js');

console.log('🧪 === DESIGN INTENT EXTRACTOR TEST SUITE ===\n');

let passCount = 0;
let failCount = 0;

function testCase(name, fn) {
  try {
    fn();
    console.log(`✅ Test: ${name}`);
    passCount++;
  } catch (error) {
    console.log(`❌ Test: ${name}`);
    console.log(`   Error: ${error.message}`);
    failCount++;
  }
}

// Test 1: 문서에서 설계 의도 추출
testCase('Test 1: Extract intent from documentation', () => {
  const extractor = new DesignIntentExtractor();

  const code = `
    /**
     * Purpose: Handle user authentication
     * Design: Use JWT tokens for stateless auth
     * Responsibility: Validate tokens and manage sessions
     * Pattern: Strategy pattern for auth providers
     */
    class AuthService {
      authenticate(user) { }
    }
  `;

  const intents = extractor.extractFromDocumentation(code, 'auth.js');

  if (intents.length < 2) {
    throw new Error(`Expected at least 2 intents, got ${intents.length}`);
  }

  const hasPurpose = intents.some(i => i.description.includes('Purpose'));
  if (!hasPurpose) {
    throw new Error('Purpose not extracted');
  }

  const hasDesign = intents.some(i => i.description.includes('Design'));
  if (!hasDesign) {
    throw new Error('Design not extracted');
  }
});

// Test 2: 이름 규칙에서 의도 파악
testCase('Test 2: Extract intent from naming conventions', () => {
  const extractor = new DesignIntentExtractor();

  const symbols = [
    'getUser',
    'setPassword',
    'isAuthenticated',
    'hasPermission',
    'canDelete',
    '_privateMethod',
    'MAX_RETRIES',
    'UserHandler',
    'SessionManager',
    'AuthController'
  ];

  const intents = extractor.extractFromNaming(symbols, 'service.js');

  if (intents.length < 8) {
    throw new Error(`Expected at least 8 intents, got ${intents.length}`);
  }

  const hasGetter = intents.some(i => i.description.includes('Getter'));
  if (!hasGetter) {
    throw new Error('Getter pattern not detected');
  }

  const hasHandler = intents.some(i => i.description.includes('Event'));
  if (!hasHandler) {
    throw new Error('Handler pattern not detected');
  }
});

// Test 3: 구조에서 패턴 의도 감지
testCase('Test 3: Extract intent from structure', () => {
  const extractor = new DesignIntentExtractor();

  const structure = {
    layers: {
      presentation: ['controller'],
      application: ['service'],
      domain: ['entity'],
      infrastructure: ['repository']
    },
    modules: ['auth', 'user', 'product', 'order'],
    dependencyDirection: 'acyclic',
    cohesion: 0.85,
    interfaces: ['IAuthService', 'IUserRepository']
  };

  const intents = extractor.extractFromStructure(structure, 'project.js');

  if (intents.length < 3) {
    throw new Error(`Expected at least 3 intents, got ${intents.length}`);
  }

  const hasLayered = intents.some(i => i.description.includes('Layered'));
  if (!hasLayered) {
    throw new Error('Layered architecture not detected');
  }

  const hasHighCohesion = intents.some(i => i.description.includes('High cohesion'));
  if (!hasHighCohesion) {
    throw new Error('High cohesion intent not detected');
  }
});

// Test 4: 구현과 의도의 괴리 감지
testCase('Test 4: Detect intent-implementation violations', () => {
  const extractor = new DesignIntentExtractor();

  const intents = [
    { description: 'High cohesion design', location: 'UserService' },
    { description: 'Private/internal: _secret', location: 'Service' }
  ];

  const implementation = {
    cohesion: 0.45, // Low cohesion - 괴리 발생!
    publicAccess: ['Service'],
    crossLayerAccess: 2
  };

  const violations = extractor.detectIntentViolation(intents, implementation);

  if (violations.length === 0) {
    throw new Error('No violations detected (should find cohesion mismatch)');
  }

  const hasHighSeverity = violations.some(v => v.severity === 'HIGH');
  if (!hasHighSeverity) {
    throw new Error('High severity violation not detected');
  }
});

// Test 5: 의도 진화 추적
testCase('Test 5: Track intent evolution', () => {
  const extractor = new DesignIntentExtractor();

  const intentHistory = [
    { description: 'Simple monolithic design', timestamp: new Date() },
    { description: 'Add separation of concerns', timestamp: new Date() },
    { description: 'Introduce service layer', timestamp: new Date() },
    { description: 'Add repository pattern', timestamp: new Date() },
    { description: 'Implement dependency injection', timestamp: new Date() },
    { description: 'Optimize for scalability', timestamp: new Date() },
    { description: 'Add caching layer', timestamp: new Date() },
    { description: 'Stabilize API contracts', timestamp: new Date() }
  ];

  const evolution = extractor.trackIntentEvolution(intentHistory);

  if (evolution.length !== intentHistory.length) {
    throw new Error(`Evolution length mismatch: ${evolution.length} vs ${intentHistory.length}`);
  }

  const initialStage = evolution[0];
  if (!initialStage.stage.includes('Initial')) {
    throw new Error(`Initial stage incorrect: ${initialStage.stage}`);
  }

  const lastStage = evolution[evolution.length - 1];
  if (!lastStage.stage.includes('Mature')) {
    throw new Error(`Final stage should be Mature, got: ${lastStage.stage}`);
  }
});

// Test 6: 설계 결정 이유 분석
testCase('Test 6: Analyze design decisions', () => {
  const extractor = new DesignIntentExtractor();

  const code = `
    class UserService {
      constructor(private repository: UserRepository) { }

      async getUser(id) {
        try {
          return await this.repository.findById(id);
        } catch (error) {
          throw new UserNotFoundError(error);
        }
      }
    }
  `;

  const patterns = ['MVC', 'Repository Pattern', 'Dependency Injection'];

  const decisions = extractor.analyzeDesignDecisions(code, patterns);

  if (decisions.length < 2) {
    throw new Error(`Expected at least 2 decisions, got ${decisions.length}`);
  }

  const hasPatternDecision = decisions.some(d => d.decision.includes('Choose'));
  if (!hasPatternDecision) {
    throw new Error('Pattern choice decision not found');
  }
});

// Test 7: 통합 설계 의도 분석
testCase('Test 7: Full integration - Design intent analysis', () => {
  const extractor = new DesignIntentExtractor();

  const code = `
    /**
     * Purpose: Manage user accounts
     * Design: Repository pattern with dependency injection
     * Responsibility: CRUD operations on users
     * Pattern: Service layer abstraction
     * TODO: Add caching for frequently accessed users
     */
    class UserService {
      constructor(private userRepository: UserRepository) { }

      getUser(id) { return this.userRepository.findById(id); }
      _internalHelper() { }
      MAX_RESULTS = 100;
    }
  `;

  const metadata = {
    filePath: 'user.service.js',
    structure: {
      layers: { application: ['service'], domain: ['entity'] },
      dependencyDirection: 'acyclic',
      cohesion: 0.88,
      interfaces: ['IUserRepository']
    },
    patterns: ['Repository', 'Dependency Injection'],
    implementation: {
      cohesion: 0.88,
      crossLayerAccess: 0
    }
  };

  const result = extractor.extractDesignIntents(code, metadata);

  if (!result.summary) {
    throw new Error('Result missing summary');
  }

  if (result.intents.length < 5) {
    throw new Error(`Expected at least 5 intents, got ${result.intents.length}`);
  }

  if (result.summary.averageConfidence < 0.75) {
    throw new Error(`Average confidence too low: ${result.summary.averageConfidence}`);
  }

  if (!result.summary.intentsByType.DOCUMENTATION) {
    throw new Error('No documentation intents');
  }
});

// Test 8: 다양한 의도 타입 분류
testCase('Test 8: Intent type grouping', () => {
  const extractor = new DesignIntentExtractor();

  const intents = [
    { type: 'DOCUMENTATION', description: 'Purpose: Auth' },
    { type: 'DOCUMENTATION', description: 'Design: JWT' },
    { type: 'NAMING', description: 'Getter: getUser' },
    { type: 'STRUCTURE', description: 'Layered architecture' },
    { type: 'NAMING', description: 'Private: _secret' }
  ];

  const grouped = extractor.groupByType(intents);

  if (grouped.DOCUMENTATION !== 2) {
    throw new Error(`DOCUMENTATION count should be 2, got ${grouped.DOCUMENTATION}`);
  }

  if (grouped.NAMING !== 2) {
    throw new Error(`NAMING count should be 2, got ${grouped.NAMING}`);
  }

  if (grouped.STRUCTURE !== 1) {
    throw new Error(`STRUCTURE count should be 1, got ${grouped.STRUCTURE}`);
  }
});

// Test 9: 심볼 추출 및 분석
testCase('Test 9: Symbol extraction and analysis', () => {
  const extractor = new DesignIntentExtractor();

  const code = `
    function authenticate(user) { }
    const validateToken = (token) => { }
    class UserService { }
    struct AuthConfig { }
  `;

  const symbols = extractor.extractSymbols(code);

  if (symbols.length < 3) {
    throw new Error(`Expected at least 3 symbols, got ${symbols.length}`);
  }

  if (!symbols.includes('authenticate')) {
    throw new Error('Function "authenticate" not extracted');
  }

  if (!symbols.includes('UserService')) {
    throw new Error('Class "UserService" not extracted');
  }
});

// Test 10: 의도맵 생성
testCase('Test 10: Intent map creation', () => {
  const extractor = new DesignIntentExtractor();

  const intents = [
    new (require('./design-intent-extractor.js')).DesignIntent('DOCUMENTATION', 'Purpose: Auth', 0.9, 'auth.js'),
    new (require('./design-intent-extractor.js')).DesignIntent('DOCUMENTATION', 'Design: JWT', 0.9, 'auth.js'),
    new (require('./design-intent-extractor.js')).DesignIntent('NAMING', 'Getter: getUser', 0.85, 'user.js'),
    new (require('./design-intent-extractor.js')).DesignIntent('STRUCTURE', 'Layered', 0.92, 'project.js')
  ];

  const map = extractor.createIntentMap(intents);

  if (!map.has('DOCUMENTATION')) {
    throw new Error('DOCUMENTATION type missing from map');
  }

  const docIntents = map.get('DOCUMENTATION');
  if (docIntents.length !== 2) {
    throw new Error(`DOCUMENTATION intents should be 2, got ${docIntents.length}`);
  }
});

// 최종 결과 출력
console.log(`\n${'='.repeat(50)}`);
console.log(`✅ Passed: ${passCount}`);
console.log(`❌ Failed: ${failCount}`);
console.log(`📊 Total: ${passCount + failCount}`);
console.log(`${'='.repeat(50)}`);

if (failCount === 0) {
  console.log('\n🎉 ALL TESTS PASSED!');
  process.exit(0);
} else {
  console.log('\n⚠️  SOME TESTS FAILED');
  process.exit(1);
}
