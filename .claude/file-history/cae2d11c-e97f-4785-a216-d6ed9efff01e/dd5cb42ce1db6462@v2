/**
 * Architecture Pattern Detector (Axis 3 Reinforcement)
 *
 * 목적: 코드 아키텍처의 패턴을 자동으로 인식하고 개선 제안
 *
 * 기능:
 * 1. MVC/MVP/MVVM 패턴 감지
 * 2. 레이어드 아키텍처 분석
 * 3. 마이크로서비스 구조 인식
 * 4. 안티 패턴 감지 (God Object, Spaghetti Code 등)
 * 5. 의존성 순환 감지
 * 6. 코드 응집도/결합도 계산
 * 7. 아키텍처 품질 점수
 * 8. 자동 리팩토링 제안
 */

class ArchitecturePattern {
  constructor(name, confidence = 0, indicators = []) {
    this.name = name;
    this.confidence = confidence;
    this.indicators = indicators;
    this.timestamp = new Date();
  }
}

class AntiPattern {
  constructor(name, severity = 'MEDIUM', location = '', suggestion = '') {
    this.name = name;
    this.severity = severity; // CRITICAL, HIGH, MEDIUM, LOW
    this.location = location;
    this.suggestion = suggestion;
    this.timestamp = new Date();
  }
}

class ArchitecturePatternDetector {
  constructor() {
    this.fileStats = {};
    this.classHierarchy = {};
    this.dependencies = {};
    this.patterns = [];
    this.antiPatterns = [];
    this.qualityScore = 0;
  }

  /**
   * Step 1: 파일 통계 분석
   */
  analyzeFileStructure(fileTree) {
    console.log('📊 Analyzing file structure...');

    const stats = {
      totalFiles: 0,
      totalLines: 0,
      avgFileSize: 0,
      fileTypeDistribution: {},
      testCoverage: 0,
      documentationRatio: 0
    };

    // 파일 유형별 분류
    let testFiles = 0;
    let sourceFiles = 0;
    let docFiles = 0;

    Object.entries(fileTree).forEach(([path, size]) => {
      stats.totalFiles++;
      stats.totalLines += size;

      const ext = path.split('.').pop();
      stats.fileTypeDistribution[ext] = (stats.fileTypeDistribution[ext] || 0) + 1;

      // 테스트 파일 감지
      if (path.includes('test') || path.includes('spec')) {
        testFiles++;
      } else if (!path.endsWith('.md') && !path.endsWith('.doc') && !path.includes('config') && !path.includes('package')) {
        sourceFiles++;
      }

      // 문서 파일 감지
      if (path.endsWith('.md') || path.endsWith('.doc')) {
        docFiles++;
      }
    });

    stats.avgFileSize = stats.totalLines / stats.totalFiles;
    // 테스트 커버리지: (테스트 파일 수 / 소스 파일 수) * 100
    stats.testCoverage = sourceFiles > 0 ? (testFiles / sourceFiles) * 100 : 0;
    stats.documentationRatio = (docFiles / stats.totalFiles) * 100;

    this.fileStats = stats;
    return stats;
  }

  /**
   * Step 2: MVC/MVP/MVVM 패턴 감지
   */
  detectMVCPattern(fileTree) {
    console.log('🎯 Detecting MVC/MVP/MVVM patterns...');

    const hasModel = Object.keys(fileTree).some(f =>
      f.includes('model') || f.includes('entity') || f.includes('schema'));

    const hasView = Object.keys(fileTree).some(f =>
      f.includes('view') || f.includes('component') || f.includes('ui'));

    const hasController = Object.keys(fileTree).some(f =>
      f.includes('controller') || f.includes('handler') || f.includes('service'));

    const hasPresenter = Object.keys(fileTree).some(f =>
      f.includes('presenter') || f.includes('presenter'));

    const hasViewModel = Object.keys(fileTree).some(f =>
      f.includes('viewmodel') || f.includes('viewstate'));

    let pattern = null;
    let confidence = 0;

    if (hasModel && hasView && hasController && !hasPresenter && !hasViewModel) {
      confidence = 0.9;
      pattern = new ArchitecturePattern('MVC', confidence, [
        'Model files detected',
        'View files detected',
        'Controller files detected'
      ]);
    } else if (hasModel && hasView && hasPresenter) {
      confidence = 0.85;
      pattern = new ArchitecturePattern('MVP', confidence, [
        'Model files detected',
        'View files detected',
        'Presenter files detected'
      ]);
    } else if (hasModel && hasViewModel) {
      confidence = 0.8;
      pattern = new ArchitecturePattern('MVVM', confidence, [
        'Model files detected',
        'ViewModel files detected'
      ]);
    }

    if (pattern) {
      this.patterns.push(pattern);
      return pattern;
    }

    return null;
  }

  /**
   * Step 3: 레이어드 아키텍처 분석
   */
  detectLayeredArchitecture(fileTree) {
    console.log('📚 Detecting layered architecture...');

    const layers = {
      presentation: [],
      application: [],
      domain: [],
      infrastructure: []
    };

    Object.keys(fileTree).forEach(file => {
      if (file.includes('presentation') || file.includes('ui') || file.includes('controller')) {
        layers.presentation.push(file);
      }
      if (file.includes('application') || file.includes('service') || file.includes('usecase')) {
        layers.application.push(file);
      }
      if (file.includes('domain') || file.includes('entity') || file.includes('model')) {
        layers.domain.push(file);
      }
      if (file.includes('infrastructure') || file.includes('repository') || file.includes('database')) {
        layers.infrastructure.push(file);
      }
    });

    const filledLayers = Object.values(layers).filter(l => l.length > 0).length;
    const confidence = Math.min(filledLayers / 4, 1.0) * 100;

    if (filledLayers >= 2) {
      const pattern = new ArchitecturePattern('Layered Architecture', confidence, [
        `Found ${filledLayers} layers`,
        `Presentation: ${layers.presentation.length} files`,
        `Application: ${layers.application.length} files`,
        `Domain: ${layers.domain.length} files`,
        `Infrastructure: ${layers.infrastructure.length} files`
      ]);
      this.patterns.push(pattern);
      return { pattern, layers };
    }

    return null;
  }

  /**
   * Step 4: 마이크로서비스 아키텍처 감지
   */
  detectMicroservicesArchitecture(fileTree) {
    console.log('🌐 Detecting microservices architecture...');

    const services = {};
    const hasServiceDirectory = Object.keys(fileTree).some(f => f.includes('services'));

    if (!hasServiceDirectory) return null;

    Object.keys(fileTree).forEach(file => {
      const match = file.match(/services\/([^/]+)\//);
      if (match) {
        services[match[1]] = (services[match[1]] || 0) + 1;
      }
    });

    const numServices = Object.keys(services).length;
    if (numServices >= 2) {
      const confidence = Math.min(numServices / 5, 1.0) * 100;
      const pattern = new ArchitecturePattern('Microservices', confidence, [
        `Found ${numServices} independent services`,
        `Services: ${Object.keys(services).join(', ')}`
      ]);
      this.patterns.push(pattern);
      return pattern;
    }

    return null;
  }

  /**
   * Step 5: 안티 패턴 감지
   */
  detectAntiPatterns(fileTree, codeQuality) {
    console.log('⚠️ Detecting anti-patterns...');

    // God Object: 파일 크기가 너무 큼
    Object.entries(fileTree).forEach(([path, size]) => {
      if (size > 1000) { // 1000줄 이상
        this.antiPatterns.push(new AntiPattern(
          'God Object',
          'HIGH',
          path,
          '파일을 작은 단위로 분리하세요 (SRP: Single Responsibility Principle)'
        ));
      }
    });

    // Spaghetti Code: 낮은 응집도, 높은 결합도
    const avgFileSize = this.fileStats.avgFileSize || 300;
    const fileVariance = Object.values(fileTree).reduce((sum, size) => {
      return sum + Math.pow(size - avgFileSize, 2);
    }, 0) / Object.keys(fileTree).length;

    if (fileVariance > 100000) {
      this.antiPatterns.push(new AntiPattern(
        'Inconsistent Module Sizes',
        'MEDIUM',
        'overall',
        '파일 크기가 일관성 없습니다. 일관된 모듈 크기를 유지하세요.'
      ));
    }

    // Circular Dependencies: 순환 의존성
    if (codeQuality && codeQuality.circularDependencies > 0) {
      this.antiPatterns.push(new AntiPattern(
        'Circular Dependencies',
        'CRITICAL',
        'dependencies',
        '순환 의존성을 제거하세요. 의존성 방향을 명확히 하세요.'
      ));
    }

    // Duplicate Code: 중복 코드
    if (codeQuality && codeQuality.duplicationRatio > 0.15) {
      this.antiPatterns.push(new AntiPattern(
        'Code Duplication',
        'HIGH',
        'overall',
        `${(codeQuality.duplicationRatio * 100).toFixed(1)}% 중복 코드. 공통 함수로 추출하세요.`
      ));
    }

    // Testing Gap: 테스트 커버리지 부족
    if (this.fileStats.testCoverage <= 30) {
      this.antiPatterns.push(new AntiPattern(
        'Low Test Coverage',
        'HIGH',
        'tests',
        `테스트 커버리지가 ${this.fileStats.testCoverage.toFixed(1)}%입니다. 80% 이상을 목표로 하세요.`
      ));
    }

    // Documentation Gap: 문서 부족
    if (this.fileStats.documentationRatio < 5) {
      this.antiPatterns.push(new AntiPattern(
        'Insufficient Documentation',
        'MEDIUM',
        'docs',
        '문서가 부족합니다. 주요 아키텍처 결정을 문서화하세요.'
      ));
    }

    return this.antiPatterns;
  }

  /**
   * Step 6: 코드 응집도 및 결합도 계산
   */
  calculateCohesionCoupling(dependencies) {
    console.log('🔗 Calculating cohesion and coupling...');

    let totalDependencies = 0;
    let internalDependencies = 0;
    let modules = Object.keys(dependencies || {}).length;

    Object.entries(dependencies || {}).forEach(([module, deps]) => {
      deps.forEach(dep => {
        totalDependencies++;
        if (dep.startsWith(module.split('/')[0])) {
          internalDependencies++;
        }
      });
    });

    const cohesion = modules > 0 ? internalDependencies / Math.max(totalDependencies, 1) : 0;
    const coupling = modules > 0 ? (totalDependencies - internalDependencies) / modules : 0;

    return {
      cohesion: Math.min(cohesion, 1.0), // 높을수록 좋음 (0-1)
      coupling: Math.min(coupling, 1.0), // 낮을수록 좋음 (0-1)
      ratio: cohesion > 0 ? cohesion / (coupling + 0.1) : 0
    };
  }

  /**
   * Step 7: 아키텍처 품질 점수 계산
   */
  calculateQualityScore() {
    console.log('⭐ Calculating architecture quality score...');

    let score = 0;
    const weights = {
      patternDetection: 0.25,
      testCoverage: 0.25,
      documentation: 0.15,
      noAntiPatterns: 0.20,
      fileConsistency: 0.15
    };

    // 1. 패턴 감지 (25%)
    const patternScore = this.patterns.length * 0.2; // 최대 5개 패턴
    score += Math.min(patternScore, 25) * weights.patternDetection / 25;

    // 2. 테스트 커버리지 (25%)
    const testScore = Math.min(this.fileStats.testCoverage, 100);
    score += (testScore / 100) * 25 * weights.testCoverage / 25;

    // 3. 문서화 (15%)
    const docScore = Math.min(this.fileStats.documentationRatio * 3, 100);
    score += (docScore / 100) * 15 * weights.documentation / 15;

    // 4. 안티 패턴 없음 (20%)
    const antiPatternPenalty = this.antiPatterns.length * 5;
    const noAntiPatternScore = Math.max(20 - antiPatternPenalty, 0);
    score += noAntiPatternScore * weights.noAntiPatterns / 20;

    // 5. 파일 일관성 (15%)
    const avgSize = this.fileStats.avgFileSize || 1;
    const fileVariance = Math.sqrt(Object.values(this.fileStats).reduce((sum, size) => {
      return sum + Math.pow((size || 0) - avgSize, 2);
    }, 0) / Math.max(this.fileStats.totalFiles, 1));

    const consistencyScore = Math.max(15 - (fileVariance / 100), 0);
    score += consistencyScore * weights.fileConsistency / 15;

    this.qualityScore = Math.round(score);
    return this.qualityScore;
  }

  /**
   * Step 8: 자동 리팩토링 제안 생성
   */
  generateRefactoringRecommendations() {
    console.log('💡 Generating refactoring recommendations...');

    const recommendations = [];

    // 우선순위: CRITICAL > HIGH > MEDIUM > LOW
    const critical = this.antiPatterns.filter(ap => ap.severity === 'CRITICAL');
    const high = this.antiPatterns.filter(ap => ap.severity === 'HIGH');
    const medium = this.antiPatterns.filter(ap => ap.severity === 'MEDIUM');

    // P0: Critical 문제 해결
    critical.forEach(ap => {
      recommendations.push({
        priority: 'P0',
        issue: ap.name,
        action: ap.suggestion,
        estimatedEffort: 'HIGH'
      });
    });

    // P1: High 문제 해결
    high.forEach(ap => {
      recommendations.push({
        priority: 'P1',
        issue: ap.name,
        action: ap.suggestion,
        estimatedEffort: 'MEDIUM'
      });
    });

    // P2: Medium 문제 해결
    medium.slice(0, 3).forEach(ap => {
      recommendations.push({
        priority: 'P2',
        issue: ap.name,
        action: ap.suggestion,
        estimatedEffort: 'LOW'
      });
    });

    // 패턴별 개선 제안
    if (this.patterns.some(p => p.name === 'MVC')) {
      recommendations.push({
        priority: 'P2',
        issue: 'MVC 패턴 최적화',
        action: 'MVC 패턴을 유지하되, Presenter 계층을 추가하여 MVP로 진화하는 것을 검토하세요.',
        estimatedEffort: 'MEDIUM'
      });
    }

    if (this.fileStats.testCoverage < 50) {
      recommendations.push({
        priority: 'P1',
        issue: '테스트 커버리지 개선',
        action: `현재 ${this.fileStats.testCoverage.toFixed(1)}% 커버리지. 80%까지 증가시키기 위해 ${Math.ceil((80 - this.fileStats.testCoverage) / 10)} 주 소요 예상.`,
        estimatedEffort: 'HIGH'
      });
    }

    return recommendations;
  }

  /**
   * 통합 분석 실행
   */
  analyzeArchitecture(fileTree, codeQuality = {}) {
    console.log('\n🏗️ === ARCHITECTURE ANALYSIS STARTED ===\n');

    const startTime = Date.now();

    // 1단계: 파일 구조 분석
    this.analyzeFileStructure(fileTree);

    // 2단계: MVC/MVP/MVVM 패턴
    this.detectMVCPattern(fileTree);

    // 3단계: 레이어드 아키텍처
    this.detectLayeredArchitecture(fileTree);

    // 4단계: 마이크로서비스 아키텍처
    this.detectMicroservicesArchitecture(fileTree);

    // 5단계: 안티 패턴 감지
    this.detectAntiPatterns(fileTree, codeQuality);

    // 6단계: 응집도/결합도 계산
    const metrics = this.calculateCohesionCoupling(codeQuality.dependencies || {});

    // 7단계: 품질 점수 계산
    this.calculateQualityScore();

    // 8단계: 리팩토링 제안
    const recommendations = this.generateRefactoringRecommendations();

    const duration = Date.now() - startTime;

    console.log('\n✅ === ARCHITECTURE ANALYSIS COMPLETED ===\n');

    return {
      summary: {
        fileStats: this.fileStats,
        patterns: this.patterns,
        antiPatterns: this.antiPatterns,
        metrics,
        qualityScore: this.qualityScore,
        recommendations
      },
      duration,
      timestamp: new Date()
    };
  }
}

module.exports = { ArchitecturePatternDetector, ArchitecturePattern, AntiPattern };
