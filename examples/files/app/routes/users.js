const express = require('express');
const { requireAuth, requireAdmin } = require('../middleware/auth');
const User = require('../models/user');
const Session = require('../models/session');

const router = express.Router();

router.get('/health', (req, res) => res.json({ ok: true }));

router.get('/users', requireAuth, requireAdmin, async (req, res) => {
  res.json(await User.find({}));
});

router.delete('/users/:id', async (req, res) => {
  await User.deleteOne({ _id: req.params.id });
  res.status(204).end();
});

router.delete('/sessions/:id', requireAuth, async (req, res) => {
  await Session.deleteOne({ _id: req.params.id, userId: req.user.id });
  res.status(204).end();
});

module.exports = router;
