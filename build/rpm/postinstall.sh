#!/usr/bin/env sh
systemctl daemon-reload
systemctl enable --now crispy.service
