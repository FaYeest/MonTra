"""Model metadata and text normalization primitives for transaction categorization."""

from __future__ import annotations

import re
from dataclasses import dataclass

_WHITESPACE_RE = re.compile(r"\s+")
_NON_ALNUM_RE = re.compile(r"[^a-z0-9\s]")


@dataclass(frozen=True)
class ModelMetadata:
    """Describes an evaluated model artifact registered with the ML service."""

    name: str
    model_version: str
    loaded: bool


def normalize_merchant_description(raw: str) -> str:
    """Normalizes a transaction description for deterministic feature extraction.

    Removes punctuation noise, lowercases ASCII/Indonesian tokens, and collapses whitespace.
    """
    lowered = raw.strip().lower()
    cleaned = _NON_ALNUM_RE.sub(" ", lowered)
    return _WHITESPACE_RE.sub(" ", cleaned).strip()
