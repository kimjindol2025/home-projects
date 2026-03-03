/**
 * Test Suite: Architecture Pattern Detector
 *
 * 8개 테스트 시나리오:
 * 1. MVC 패턴 감지
 * 2. 레이어드 아키텍처 감지
 * 3. 마이크로서비스 감지
 * 4. God Object 안티 패턴 감지
 * 5. Circular Dependencies 감지
 * 6. 낮은 테스트 커버리지 감지
 * 7. 품질 점수 계산
 * 8. 자동 리팩토링 제안 생성
 */

const { ArchitecturePatternDetector } = require('./architecture-pattern-detector.js');

console.log('🧪 === ARCHITECTURE DETECTOR TEST SUITE ===\n');

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

// Test 1: MVC 패턴 감지
testCase('Test 1: MVC Pattern Detection', () => {
  const detector = new ArchitecturePatternDetector();

  const mvcFileTree = {
    'src/models/User.js': 150,
    'src/models/Product.js': 120,
    'src/views/UserView.js': 200,
    'src/views/ProductView.js': 180,
    'src/controllers/UserController.js': 300,
    'src/controllers/ProductController.js': 250,
    'src/routes/api.js': 80
  };

  const pattern = detector.detectMVCPattern(mvcFileTree);

  if (!pattern || pattern.name !== 'MVC') {
    throw new Error('Failed to detect MVC pattern');
  }
  if (pattern.confidence < 0.8) {
    throw new Error(`MVC confidence too low: ${pattern.confidence}`);
  }
});

// Test 2: 레이어드 아키텍처 감지
testCase('Test 2: Layered Architecture Detection', () => {
  const detector = new ArchitecturePatternDetector();

  const layeredFileTree = {
    'src/presentation/UserController.js': 200,
    'src/application/UserService.js': 250,
    'src/domain/User.js': 150,
    'src/infrastructure/UserRepository.js': 180,
    'src/database/schema.sql': 100
  };

  const result = detector.detectLayeredArchitecture(layeredFileTree);

  if (!result || result.pattern.name !== 'Layered Architecture') {
    throw new Error('Failed to detect layered architecture');
  }
  if (result.layers.presentation.length === 0) {
    throw new Error('Presentation layer not detected');
  }
});

// Test 3: 마이크로서비스 아키텍처 감지
testCase('Test 3: Microservices Architecture Detection', () => {
  const detector = new ArchitecturePatternDetector();

  const microservicesFileTree = {
    'services/user-service/index.js': 300,
    'services/user-service/db.js': 150,
    'services/user-service/tests/unit.js': 200,
    'services/product-service/index.js': 280,
    'services/product-service/db.js': 140,
    'services/product-service/tests/unit.js': 180,
    'services/order-service/index.js': 320,
    'shared/utils.js': 100
  };

  const pattern = detector.detectMicroservicesArchitecture(microservicesFileTree);

  if (!pattern || pattern.name !== 'Microservices') {
    throw new Error('Failed to detect microservices architecture');
  }
  if (pattern.confidence < 0.6) {
    throw new Error(`Microservices confidence too low: ${pattern.confidence}`);
  }
});

// Test 4: God Object 안티 패턴 감지
testCase('Test 4: God Object Anti-pattern Detection', () => {
  const detector = new ArchitecturePatternDetector();

  const fileTree = {
    'src/GodObject.js': 5000, // 5000줄 - 너무 큼
    'src/User.js': 150,
    'src/Product.js': 120,
    'tests/unit.js': 100
  };

  detector.analyzeFileStructure(fileTree);
  detector.detectAntiPatterns(fileTree, {});

  const godObjects = detector.antiPatterns.filter(ap => ap.name === 'God Object');

  if (godObjects.length === 0) {
    throw new Error('God Object anti-pattern not detected');
  }
  if (godObjects[0].severity !== 'HIGH') {
    throw new Error('God Object severity should be HIGH');
  }
});

// Test 5: 순환 의존성 감지
testCase('Test 5: Circular Dependencies Detection', () => {
  const detector = new ArchitecturePatternDetector();

  const fileTree = {
    'src/A.js': 100,
    'src/B.js': 100,
    'src/C.js': 100
  };

  const codeQuality = {
    circularDependencies: 2
  };

  detector.analyzeFileStructure(fileTree);
  detector.detectAntiPatterns(fileTree, codeQuality);

  const circularDeps = detector.antiPatterns.filter(ap => ap.name === 'Circular Dependencies');

  if (circularDeps.length === 0) {
    throw new Error('Circular dependencies not detected');
  }
  if (circularDeps[0].severity !== 'CRITICAL') {
    throw new Error('Circular dependencies should be CRITICAL');
  }
});

// Test 6: 낮은 테스트 커버리지 감지
testCase('Test 6: Low Test Coverage Detection', () => {
  const detector = new ArchitecturePatternDetector();

  const fileTree = {
    'src/User.js': 150,
    'src/Product.js': 120,
    'src/Order.js': 200,
    'src/Payment.js': 180,
    'src/Shipping.js': 100,
    'tests/user.test.js': 100  // 1/5 = 20% coverage
  };

  detector.analyzeFileStructure(fileTree);
  console.log(`Test Coverage: ${detector.fileStats.testCoverage.toFixed(1)}%`);
  detector.detectAntiPatterns(fileTree, {});

  const testGaps = detector.antiPatterns.filter(ap => ap.name === 'Low Test Coverage');

  if (testGaps.length === 0) {
    throw new Error(`Low test coverage not detected. Coverage: ${detector.fileStats.testCoverage.toFixed(1)}%`);
  }
  if (testGaps[0].severity !== 'HIGH') {
    throw new Error('Low test coverage should be HIGH');
  }
});

// Test 7: 품질 점수 계산
testCase('Test 7: Quality Score Calculation', () => {
  const detector = new ArchitecturePatternDetector();

  const fileTree = {
    'src/models/User.js': 150,
    'src/controllers/UserController.js': 200,
    'src/views/UserView.js': 180,
    'tests/unit.js': 200,
    'docs/architecture.md': 50
  };

  detector.analyzeFileStructure(fileTree);
  detector.detectMVCPattern(fileTree);
  detector.detectAntiPatterns(fileTree, {});
  const score = detector.calculateQualityScore();

  if (score < 0 || score > 100) {
    throw new Error(`Quality score out of range: ${score}`);
  }
  if (score < 50) {
    throw new Error(`Quality score too low for well-structured project: ${score}`);
  }
});

// Test 8: 자동 리팩토링 제안 생성
testCase('Test 8: Automatic Refactoring Recommendations', () => {
  const detector = new ArchitecturePatternDetector();

  const fileTree = {
    'src/BigFile.js': 2000, // God Object
    'src/util.js': 800,
    'tests/test.js': 50 // 낮은 테스트 커버리지
  };

  detector.analyzeFileStructure(fileTree);
  detector.detectAntiPatterns(fileTree, {});
  const recommendations = detector.generateRefactoringRecommendations();

  if (recommendations.length === 0) {
    throw new Error('No refactoring recommendations generated');
  }

  const hasPriority = recommendations.some(r => r.priority === 'P0' || r.priority === 'P1');
  if (!hasPriority) {
    throw new Error('No high-priority recommendations found');
  }
});

// 통합 테스트: 실제 프로젝트 시뮬레이션
testCase('Test 9: Full Integration - Complex Project Analysis', () => {
  const detector = new ArchitecturePatternDetector();

  const complexProjectFileTree = {
    // Presentation Layer
    'src/presentation/UserController.js': 250,
    'src/presentation/ProductController.js': 220,
    'src/presentation/views/UserList.js': 180,
    'src/presentation/views/ProductList.js': 160,

    // Application Layer
    'src/application/UserService.js': 300,
    'src/application/ProductService.js': 280,
    'src/application/dtos/UserDTO.js': 100,

    // Domain Layer
    'src/domain/User.js': 200,
    'src/domain/Product.js': 180,
    'src/domain/Order.js': 150,

    // Infrastructure Layer
    'src/infrastructure/UserRepository.js': 220,
    'src/infrastructure/ProductRepository.js': 200,
    'src/infrastructure/database/connection.js': 100,

    // Tests
    'tests/unit/UserService.test.js': 300,
    'tests/unit/ProductService.test.js': 280,
    'tests/integration/api.test.js': 400,

    // Documentation
    'docs/architecture.md': 150,
    'docs/api.md': 120,
    'docs/database.md': 100,

    // Config
    'config/database.js': 50,
    'package.json': 30
  };

  const codeQuality = {
    duplicationRatio: 0.08,
    circularDependencies: 0,
    dependencies: {
      'src/presentation': ['src/application'],
      'src/application': ['src/domain', 'src/infrastructure'],
      'src/domain': [],
      'src/infrastructure': ['src/domain']
    }
  };

  const result = detector.analyzeArchitecture(complexProjectFileTree, codeQuality);

  // 기본 검증
  if (!result.summary) {
    throw new Error('Analysis result missing summary');
  }
  if (result.summary.qualityScore < 60) {
    throw new Error(`Quality score too low for well-structured project: ${result.summary.qualityScore}`);
  }
  if (result.summary.patterns.length === 0) {
    throw new Error('No patterns detected in well-structured project');
  }
});

// Test 10: 응집도/결합도 계산
testCase('Test 10: Cohesion and Coupling Metrics', () => {
  const detector = new ArchitecturePatternDetector();

  const dependencies = {
    'moduleA': ['moduleB', 'external'],
    'moduleB': ['moduleC'],
    'moduleC': []
  };

  const metrics = detector.calculateCohesionCoupling(dependencies);

  if (metrics.cohesion < 0 || metrics.cohesion > 1) {
    throw new Error(`Cohesion out of range: ${metrics.cohesion}`);
  }
  if (metrics.coupling < 0 || metrics.coupling > 1) {
    throw new Error(`Coupling out of range: ${metrics.coupling}`);
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
