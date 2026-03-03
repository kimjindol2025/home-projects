/**
 * Design Intent Extractor (Phase C: Axis 3 강화)
 *
 * 목적: 코드에서 설계 의도를 자동으로 추출하고 분석
 *
 * 기능:
 * 1. 주석/문서에서 설계 의도 추출
 * 2. 이름 규칙에서 의도 파악
 * 3. 구조에서 패턴 의도 감지
 * 4. 구현과 의도의 괴리 감지
 * 5. 의도 진화 추적
 * 6. 설계 결정 이유 분석
 */

class DesignIntent {
  constructor(type, description, confidence, location) {
    this.type = type; // NAMING, STRUCTURE, DOCUMENTATION, PATTERN
    this.description = description;
    this.confidence = confidence; // 0-1
    this.location = location;
    this.timestamp = new Date();
  }
}

class DesignDecision {
  constructor(decision, rationale, alternatives, chosen) {
    this.decision = decision;
    this.rationale = rationale;
    this.alternatives = alternatives;
    this.chosen = chosen;
    this.timestamp = new Date();
  }
}

class DesignIntentExtractor {
  constructor() {
    this.intents = [];
    this.decisions = [];
    this.intentMap = new Map();
    this.intentEvolution = [];
  }

  /**
   * Step 1: 주석/문서에서 설계 의도 추출
   */
  extractFromDocumentation(code, filePath) {
    console.log('📝 Extracting design intent from documentation...');

    const intents = [];

    // 패턴 1: "목적:" 또는 "Purpose:"
    const purposeMatch = code.match(/[Pp]urpose[:\s]+([^\n]+)/);
    if (purposeMatch) {
      intents.push(new DesignIntent(
        'DOCUMENTATION',
        `Purpose: ${purposeMatch[1].trim()}`,
        0.95,
        filePath
      ));
    }

    // 패턴 2: "설계:" 또는 "Design:"
    const designMatch = code.match(/[Dd]esign[:\s]+([^\n]+)/);
    if (designMatch) {
      intents.push(new DesignIntent(
        'DOCUMENTATION',
        `Design: ${designMatch[1].trim()}`,
        0.95,
        filePath
      ));
    }

    // 패턴 3: "책임:" 또는 "Responsibility:"
    const respMatch = code.match(/[Rr]esponsibilit[y|ies][:\s]+([^\n]+)/);
    if (respMatch) {
      intents.push(new DesignIntent(
        'DOCUMENTATION',
        `Responsibility: ${respMatch[1].trim()}`,
        0.90,
        filePath
      ));
    }

    // 패턴 4: "패턴:" 또는 "Pattern:"
    const patternMatch = code.match(/[Pp]attern[:\s]+([^\n]+)/);
    if (patternMatch) {
      intents.push(new DesignIntent(
        'DOCUMENTATION',
        `Pattern: ${patternMatch[1].trim()}`,
        0.85,
        filePath
      ));
    }

    // 패턴 5: TODO/FIXME/NOTE
    const todoMatches = code.matchAll(/(?:TODO|FIXME|NOTE)[:\s]*([^\n]+)/g);
    for (const match of todoMatches) {
      intents.push(new DesignIntent(
        'DOCUMENTATION',
        `Future Intent: ${match[1].trim()}`,
        0.70,
        filePath
      ));
    }

    return intents;
  }

  /**
   * Step 2: 이름 규칙에서 의도 파악
   */
  extractFromNaming(symbols, filePath) {
    console.log('📛 Extracting design intent from naming conventions...');

    const intents = [];

    symbols.forEach(symbol => {
      // 패턴 1: 접두사 분석
      if (symbol.startsWith('get')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Getter pattern for: ${symbol}`,
          0.85,
          filePath
        ));
      }
      if (symbol.startsWith('set')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Setter pattern for: ${symbol}`,
          0.85,
          filePath
        ));
      }
      if (symbol.startsWith('is')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Boolean predicate: ${symbol}`,
          0.90,
          filePath
        ));
      }
      if (symbol.startsWith('has')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Possession check: ${symbol}`,
          0.88,
          filePath
        ));
      }
      if (symbol.startsWith('can')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Capability check: ${symbol}`,
          0.85,
          filePath
        ));
      }

      // 패턴 2: 언더스코어 (private)
      if (symbol.startsWith('_')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Private/internal: ${symbol}`,
          0.92,
          filePath
        ));
      }

      // 패턴 3: 대문자 (상수)
      if (/^[A-Z_]+$/.test(symbol)) {
        intents.push(new DesignIntent(
          'NAMING',
          `Constant/config: ${symbol}`,
          0.90,
          filePath
        ));
      }

      // 패턴 4: Handler/Manager/Controller
      if (symbol.includes('Handler')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Event/Request handler: ${symbol}`,
          0.88,
          filePath
        ));
      }
      if (symbol.includes('Manager')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Resource management: ${symbol}`,
          0.85,
          filePath
        ));
      }
      if (symbol.includes('Controller')) {
        intents.push(new DesignIntent(
          'NAMING',
          `Flow control: ${symbol}`,
          0.87,
          filePath
        ));
      }
    });

    return intents;
  }

  /**
   * Step 3: 구조에서 패턴 의도 감지
   */
  extractFromStructure(structure, filePath) {
    console.log('🏛️ Extracting design intent from structure...');

    const intents = [];

    // 패턴 1: 계층적 구조
    const hasLayers = structure.layers && Object.keys(structure.layers).length > 0;
    if (hasLayers) {
      const layers = Object.keys(structure.layers);
      intents.push(new DesignIntent(
        'STRUCTURE',
        `Layered architecture: ${layers.join(' → ')}`,
        0.92,
        filePath
      ));
    }

    // 패턴 2: 모듈 분리
    const hasModules = structure.modules && structure.modules.length > 2;
    if (hasModules) {
      intents.push(new DesignIntent(
        'STRUCTURE',
        `Module separation: ${structure.modules.length} independent modules`,
        0.85,
        filePath
      ));
    }

    // 패턴 3: 의존성 방향
    if (structure.dependencyDirection === 'acyclic') {
      intents.push(new DesignIntent(
        'STRUCTURE',
        `Acyclic dependency graph (DAG)`,
        0.95,
        filePath
      ));
    }

    // 패턴 4: 응집도
    if (structure.cohesion > 0.8) {
      intents.push(new DesignIntent(
        'STRUCTURE',
        `High cohesion design (${(structure.cohesion * 100).toFixed(1)}%)`,
        0.90,
        filePath
      ));
    }

    // 패턴 5: 인터페이스 기반
    if (structure.interfaces && structure.interfaces.length > 0) {
      intents.push(new DesignIntent(
        'STRUCTURE',
        `Contract-based design with interfaces`,
        0.88,
        filePath
      ));
    }

    return intents;
  }

  /**
   * Step 4: 구현과 의도의 괴리 감지
   */
  detectIntentViolation(intents, actualImplementation) {
    console.log('⚠️ Detecting intent-implementation mismatches...');

    const violations = [];

    intents.forEach(intent => {
      // Intent: "High cohesion" but actual cohesion is low
      if (intent.description.includes('High cohesion') &&
          actualImplementation.cohesion < 0.6) {
        violations.push({
          intent: intent.description,
          reality: `Actual cohesion: ${(actualImplementation.cohesion * 100).toFixed(1)}%`,
          severity: 'HIGH',
          gap: 0.8 - actualImplementation.cohesion
        });
      }

      // Intent: "Private" but symbol is accessed from outside
      if (intent.description.includes('Private') &&
          actualImplementation.publicAccess &&
          actualImplementation.publicAccess.includes(intent.location)) {
        violations.push({
          intent: intent.description,
          reality: `Actually accessed publicly`,
          severity: 'HIGH',
          gap: 1.0
        });
      }

      // Intent: "Layer separation" but cross-layer access exists
      if (intent.description.includes('Layered') &&
          actualImplementation.crossLayerAccess > 0) {
        violations.push({
          intent: intent.description,
          reality: `${actualImplementation.crossLayerAccess} cross-layer accesses`,
          severity: 'MEDIUM',
          gap: 0.5
        });
      }
    });

    return violations;
  }

  /**
   * Step 5: 의도 진화 추적
   */
  trackIntentEvolution(intentHistory) {
    console.log('📈 Tracking design intent evolution...');

    const evolution = [];

    // 초기 의도 → 진화된 의도
    const stages = [
      { stage: 'Initial', description: 'First design decision' },
      { stage: 'Expanding', description: 'Adding new features' },
      { stage: 'Refining', description: 'Improving existing design' },
      { stage: 'Specializing', description: 'Optimizing for specific cases' },
      { stage: 'Mature', description: 'Stable, well-established design' }
    ];

    let currentStage = 0;
    intentHistory.forEach((intent, index) => {
      // 2-3개씩 단계별로 증가
      const stageIndex = Math.floor(index / (intentHistory.length / stages.length));
      currentStage = Math.min(stageIndex, stages.length - 1);

      evolution.push({
        version: index + 1,
        intent: intent.description,
        stage: stages[currentStage].stage,
        stageDescription: stages[currentStage].description,
        timestamp: intent.timestamp
      });
    });

    return evolution;
  }

  /**
   * Step 6: 설계 결정 이유 분석
   */
  analyzeDesignDecisions(code, patterns) {
    console.log('🎯 Analyzing design decisions...');

    const decisions = [];

    // 결정 1: 패턴 선택
    patterns.forEach(pattern => {
      decisions.push(new DesignDecision(
        `Choose ${pattern}`,
        `Provides clear separation of concerns and maintainability`,
        ['Alternative A', 'Alternative B'],
        pattern
      ));
    });

    // 결정 2: 계층 선택
    if (code.includes('presentation') || code.includes('controller')) {
      decisions.push(new DesignDecision(
        `Use Layered Architecture`,
        `Improves scalability and testability`,
        ['Microservices', 'Monolithic'],
        'Layered'
      ));
    }

    // 결정 3: 데이터 흐름
    if (code.includes('service') && code.includes('repository')) {
      decisions.push(new DesignDecision(
        `Use Dependency Injection`,
        `Reduces coupling and improves testability`,
        ['Locator pattern', 'Hard coupling'],
        'Dependency Injection'
      ));
    }

    // 결정 4: 에러 처리
    if (code.includes('try') || code.includes('catch')) {
      decisions.push(new DesignDecision(
        `Exception-based error handling`,
        `Clear error propagation and handling`,
        ['Result type', 'Error codes'],
        'Exceptions'
      ));
    }

    return decisions;
  }

  /**
   * 통합 추출 및 분석
   */
  extractDesignIntents(code, metadata = {}) {
    console.log('\n🎨 === DESIGN INTENT EXTRACTION STARTED ===\n');

    const startTime = Date.now();

    // Step 1: 문서에서 추출
    const docIntents = this.extractFromDocumentation(code, metadata.filePath || 'unknown');
    this.intents.push(...docIntents);

    // Step 2: 이름 규칙에서 추출
    const symbols = metadata.symbols || this.extractSymbols(code);
    const namingIntents = this.extractFromNaming(symbols, metadata.filePath || 'unknown');
    this.intents.push(...namingIntents);

    // Step 3: 구조에서 추출
    const structure = metadata.structure || {};
    const structureIntents = this.extractFromStructure(structure, metadata.filePath || 'unknown');
    this.intents.push(...structureIntents);

    // Step 4: 괴리 감지
    const violations = this.detectIntentViolation(
      this.intents,
      metadata.implementation || {}
    );

    // Step 5: 진화 추적
    const evolution = this.trackIntentEvolution(this.intents);
    this.intentEvolution = evolution;

    // Step 6: 결정 분석
    const decisions = this.analyzeDesignDecisions(code, metadata.patterns || []);
    this.decisions = decisions;

    // 통합 맵 생성
    this.intentMap = this.createIntentMap(this.intents);

    const duration = Date.now() - startTime;

    console.log('\n✅ === DESIGN INTENT EXTRACTION COMPLETED ===\n');

    return {
      intents: this.intents,
      violations,
      evolution,
      decisions,
      intentMap: Array.from(this.intentMap.entries()),
      summary: {
        totalIntents: this.intents.length,
        intentsByType: this.groupByType(this.intents),
        violationCount: violations.length,
        criticality: violations.filter(v => v.severity === 'HIGH').length,
        averageConfidence: this.intents.length > 0
          ? this.intents.reduce((sum, i) => sum + i.confidence, 0) / this.intents.length
          : 0
      },
      duration,
      timestamp: new Date()
    };
  }

  /**
   * 헬퍼: 심볼 추출
   */
  extractSymbols(code) {
    const symbols = [];

    // 함수/메서드명 추출
    const fnMatches = code.matchAll(/(?:function|fn|const|let)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g);
    for (const match of fnMatches) {
      symbols.push(match[1]);
    }

    // 클래스명 추출
    const classMatches = code.matchAll(/(?:class|struct)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g);
    for (const match of classMatches) {
      symbols.push(match[1]);
    }

    return [...new Set(symbols)];
  }

  /**
   * 헬퍼: 의도맵 생성
   */
  createIntentMap(intents) {
    const map = new Map();

    intents.forEach(intent => {
      if (!map.has(intent.type)) {
        map.set(intent.type, []);
      }
      map.get(intent.type).push(intent.description);
    });

    return map;
  }

  /**
   * 헬퍼: 타입별 분류
   */
  groupByType(intents) {
    const grouped = {};

    intents.forEach(intent => {
      if (!grouped[intent.type]) {
        grouped[intent.type] = 0;
      }
      grouped[intent.type]++;
    });

    return grouped;
  }
}

module.exports = {
  DesignIntentExtractor,
  DesignIntent,
  DesignDecision
};
