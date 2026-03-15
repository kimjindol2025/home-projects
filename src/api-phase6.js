/**
 * GOGS Architect Phase 6 API
 * 자기호스팅 상태 엔드포인트
 */

const express = require('express');
const router = express.Router();
const { SelfHostingBridge } = require('./self-hosting-bridge');

const bridge = new SelfHostingBridge(
  process.env.FREELANG_PATH || '/tmp/freelang-v2',
  process.env.GOGS_PATH || '/tmp/gogs-architect'
);

/**
 * GET /api/v2/self-hosting/status
 * 자기호스팅 상태 조회
 */
router.get('/self-hosting/status', async (req, res) => {
  try {
    const status = await bridge.reportToGogs();
    res.json({
      status: 'success',
      data: status,
    });
  } catch (error) {
    res.status(500).json({
      status: 'error',
      error: error.message,
    });
  }
});

/**
 * POST /api/v2/self-hosting/validate
 * 자기호스팅 검증 실행
 */
router.post('/self-hosting/validate', async (req, res) => {
  try {
    const validation = await bridge.validateSelfHosting();
    res.json({
      status: 'success',
      data: validation,
    });
  } catch (error) {
    res.status(500).json({
      status: 'error',
      error: error.message,
    });
  }
});

module.exports = router;
