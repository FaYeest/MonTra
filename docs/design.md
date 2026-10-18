# MonTra — Android App Design Specification

**Document type:** Product design specification  
**Version:** 1.0  
**Status:** Design baseline for implementation  
**Product:** MonTra — Money Tracker  
**Primary platform:** Android  
**Frontend implementation:** Flutter / Dart  
**Related document:** [`PRD.md`](./MonTra_PRD.md)  
**Prototype reference:** [`montra-android-preview.html`](./montra-android-preview.html)

---

## 1. Purpose

This document defines the visual language, interaction rules, screen structure, accessibility baseline, and acceptance criteria for the MonTra Android application. It is intended to keep design and implementation consistent as the HTML prototype is translated into Flutter.

The current HTML preview is a **visual and interaction reference, not a production implementation**. Its sample figures are illustrative and are not connected to a live database or ML model. Flutter screens must use backend data and real loading, empty, error, and success states when integrated.

For “10 golden rules of UI/UX,” this specification uses **Jakob Nielsen’s 10 usability heuristics** as the review framework. They are applied below to MonTra’s actual financial workflows, rather than treated as generic slogans.

## 2. Product design direction

### 2.1 Design intent

MonTra should feel like a dependable personal finance tool: calm, clear, precise, and human. The interface should help people understand their money without making financial tracking feel like accounting software.

The interface prioritizes:

1. Financial numbers that are easy to scan and compare.
2. Fast transaction entry with minimal friction.
3. Clear budget status and actionable explanations.
4. AI/ML output that can be checked against underlying data.
5. Privacy, control, and predictable behavior.

### 2.2 Visual character

Use a warm off-white canvas, dark forest-green anchors, restrained mint surfaces, and a small lime accent for primary actions. Typography is compact but legible, with a clear distinction between headings, supporting labels, and financial values. Use borders and spacing before using shadows to separate content.

The interface should feel designed for a real financial product, not a landing-page template pasted into a mobile app.

### 2.3 Explicit anti-patterns

Do not introduce the following unless user research demonstrates a real need:

- Neon gradients, purple-blue “AI” gradients, blurred glass panels, or glowing borders.
- Decorative 3D illustrations, random floating shapes, and excessive sparkle icons.
- A chatbot that occupies the home screen or blocks common tasks.
- A card for every single piece of text, with no meaningful hierarchy.
- Large empty hero areas inside the app that push transactions below the fold.
- Charts without a title, period, unit, legend, or interpretable meaning.
- Animated numbers that make balances hard to read.
- Generic motivational copy such as “You’re crushing it!” when neutral facts would be more useful.
- Fake AI labels, invented confidence scores, unsupported predictions, or insights with no visible evidence.
- Color as the only way to communicate income, expense, warning, or success.

## 3. Users and core jobs

The first release is optimized for students, early-career workers, and people who want more control over day-to-day spending.

Primary user jobs:

- Record an income or expense quickly after it happens.
- See how much money is recorded across accounts.
- Understand income and expense totals for a selected period.
- Find a transaction using text search and filters.
- See how much of each budget has been used and how much remains.
- Understand why an insight or warning was shown.
- Hide sensitive balances from people nearby.

The default experience should work without using AI. Users must always be able to record, edit, search, and categorize transactions manually.

## 4. Information architecture

### 4.1 Primary navigation

Use a persistent bottom navigation bar with five destinations:

1. **Beranda** — financial overview and recent transactions.
2. **Transaksi** — searchable transaction history.
3. **Tambah** — the central, visually emphasized action to record a transaction.
4. **Anggaran** — overall and category budgets.
5. **Insight** — trends, forecasts, and explanations.

The user profile and app preferences are opened through the avatar in the top app bar. Profile should not become a sixth bottom-navigation item.

### 4.2 Navigation behavior

- Preserve the selected destination state while moving between primary sections where practical.
- Reset a screen to the top only when entering that screen is expected to start a fresh browse; do not unexpectedly discard search text or filters when the user returns.
- All navigation items require a visible label and an accessible semantic label.
- The selected navigation destination must be indicated by more than a subtle color change alone. Use color plus weight, indicator, or another clear visual state.
- The center Add action opens the transaction form from any primary destination.
- Pressing Android Back closes the keyboard first, then a modal or sheet, then navigates back according to the current route. It must not silently discard a completed transaction.

## 5. Design tokens

Centralize tokens in Flutter theme and shared design-system files. Do not define separate, slightly different values inside every widget.

### 5.1 Color palette — light theme

| Token | Hex | Intended use |
|---|---|---|
| `canvas` | `#F4F3ED` | Main app background |
| `surface` | `#FFFEFA` | Cards, sheets, navigation surface |
| `surfaceSubtle` | `#EEEDE5` | Secondary surfaces, icon backgrounds, inactive tracks |
| `textPrimary` | `#202B25` | Main text and numbers |
| `textSecondary` | `#747B73` | Supporting labels and metadata |
| `outline` | `#DEDFD5` | Dividers and subtle outlines |
| `primary` | `#244D3D` | Primary actions, active navigation, key surfaces |
| `primaryAlt` | `#315F49` | Secondary emphasis when needed |
| `primaryContainer` | `#DCE9D7` | Soft green surfaces |
| `accent` | `#D7F17E` | High-priority add action and limited emphasis |
| `expense` | `#B95745` | Expense/warning text and indicators when appropriate |
| `expenseContainer` | `#F6E4DF` | Soft expense surface |
| `warning` | `#956D2D` | Budget warnings and caution |

These values are the prototype’s starting palette, not proof of sufficient contrast. Verify text, icon, chart, and control contrast in the actual app before release. Adjust tokens when necessary without changing the overall character.

### 5.2 Color palette — dark theme

| Token | Hex | Intended use |
|---|---|---|
| `canvas` | `#151B18` | Main app background |
| `surface` | `#1D2621` | Cards, sheets, navigation surface |
| `surfaceSubtle` | `#26312A` | Secondary surfaces |
| `textPrimary` | `#F2F1E9` | Main text and numbers |
| `textSecondary` | `#A3AEA4` | Supporting labels |
| `outline` | `#354139` | Dividers and subtle outlines |
| `primary` | `#D7F17E` | Primary actions and active states |
| `primaryContainer` | `#32483A` | Soft green surfaces |
| `expense` | `#E58A76` | Expense indicators |
| `warning` | `#E2B875` | Budget warnings |

Dark mode is a deliberate color mapping, not a simple inversion. Keep financial values and chart data legible without using full white backgrounds or neon effects.

### 5.3 Semantic color rules

- Green is not a universal synonym for “good”; use it for primary identity and positive movement only when the text explains the meaning.
- Income and expenses must have a sign or label as well as color, such as `+Rp500.000` and `−Rp25.000`.
- Warnings must include meaningful text such as “Mendekati batas,” not only an amber bar.
- Never imply a forecast is certain through a success-green badge.
- Verify contrast independently in both themes.

### 5.4 Typography

Use the Android system sans-serif or a bundled Inter family if it is explicitly included and licensed. Do not mix multiple typefaces for decoration.

| Style | Recommended size | Weight | Use |
|---|---:|---:|---|
| Screen title | 26–30 sp | 700–750 | Main page heading |
| Section title | 18–20 sp | 650–750 | Section heading |
| Large financial value | 28–34 sp | 700–800 | Primary balance |
| Secondary financial value | 18–22 sp | 650–750 | Income, expense, and summary values |
| Body | 14–16 sp | 400–500 | Explanations and supporting text |
| Supporting text | 12–14 sp | 400–500 | Transaction metadata and secondary labels |
| Compact labels | 12 sp minimum by default | 500–700 | Period, status, filter, chart labels |

Some type in the HTML prototype is intentionally very small to fit the preview frame. Do **not** copy its 9–11 px text sizes directly into the Android application. Use scalable `sp`, test Android font scaling, and allow long Indonesian labels to wrap or truncate deliberately.

Use tabular numerals for money values where supported. Keep currency formatting consistent, for example `Rp2.840.000`; do not mix `2.84M`, `Rp 2,840,000`, and `Rp2.840.000` on the same screen.

### 5.5 Spacing scale

Use a 4 dp base scale:

`4, 8, 12, 16, 20, 24, 32, 40, 48`

Default page horizontal padding: 16–20 dp, adjusted for narrow devices and larger font settings. Use 8–12 dp for related elements, 16–24 dp between content groups, and 24–32 dp before a distinct section.

### 5.6 Shapes, borders, and elevation

- Small controls: 10–12 dp radius.
- List and compact summary surfaces: 12–16 dp radius.
- Main balance surface or prominent insight: 18–20 dp radius.
- Bottom sheets: 20–24 dp radius on the top corners.
- Prefer a 1 dp outline or surface contrast to excessive shadows.
- Use elevation only to explain layering, such as a modal sheet or floating add action.
- Use one consistent outline approach; do not combine heavy borders and heavy shadows on the same surface.

### 5.7 Icons and motion

- Use one coherent outline icon family, ideally Material Symbols or an equivalent consistent set.
- Icons should clarify meaning. Do not add icons solely to decorate every heading.
- Icon-only buttons require accessible labels and an adequately sized touch region.
- Use motion for transitions and confirmation, not as decoration. Prefer short, restrained transitions and respect Android’s reduced-motion settings where available.
- Avoid animated balance counting, bouncing cards, parallax, and continuously moving charts.

## 6. Screen specifications

### 6.1 Home — Beranda

**Primary purpose:** Let users understand their current financial position quickly and reach common actions.

Content order:

1. Top app bar with MonTra mark, theme action if retained, and avatar/profile access.
2. Greeting and current reporting period selector.
3. “Total saldo tercatat” balance surface with a hide/show control.
4. Two compact summary values: Pemasukan and Pengeluaran.
5. Spending trend for the selected period, with a clear comparison to the prior period.
6. Recent transactions with merchant/description, category, time/date, sign, and amount.
7. Link to the complete transaction list and to relevant insights.

Behavior:

- Hide-balance state masks all prominent monetary values on the current surface and should remain consistent until the user shows them again or the app’s privacy policy dictates otherwise.
- Period selection updates all period-dependent metrics consistently.
- Tapping a recent transaction opens its details when implemented.
- The chart has an accessible description and should not be the only place a number is communicated.
- Do not use an arbitrary positive percentage without a period and comparison context.

### 6.2 Transactions — Transaksi

**Primary purpose:** Find, verify, and manage recorded transactions.

Content order:

1. Screen title and one-sentence purpose.
2. Current-period income and expense summary.
3. Search field.
4. Filter chips: Semua, Pengeluaran, Pemasukan. Add more filters only when supported by real user needs.
5. Date-grouped transaction list.
6. Empty state when no transaction matches.
7. Access to transaction details/editing when available.

Transaction row content:

- A familiar merchant/category icon if useful.
- Merchant or description, visible without unnecessary truncation.
- Category and date/time metadata.
- Signed amount aligned consistently at the end of the row.

Behavior:

- Search filters as the user types and supports a clear action.
- Active filters must be visibly selected and announced accessibly.
- Empty search results must distinguish “no transactions yet” from “no matching transactions.”
- If the user edits or deletes a transaction, update affected totals only after the backend confirms the change.
- Destructive actions require a clear consequence and an undo path where possible.

### 6.3 Add transaction — Tambah

**Primary purpose:** Record a transaction with minimal effort and prevent common data-entry mistakes.

Open as a bottom sheet or a dedicated full-screen form when the viewport or keyboard makes a sheet cramped.

Fields:

1. Transaction type: Pengeluaran or Pemasukan.
2. Amount, with numeric keyboard and clear currency context.
3. Description or merchant.
4. Category, with a sensible default but no hidden forced categorization.
5. Date, defaulted to the current local date.
6. Account/source, required when multiple accounts are implemented.
7. Optional notes only if they add real value.

Primary action: **Simpan transaksi**.

Behavior and validation:

- Put focus on the amount field when the form opens, unless restoring an in-progress draft.
- Amount must be greater than zero and valid for the selected currency rules.
- Description must be validated before saving.
- Switch the category suggestions when transaction type changes, without unexpectedly deleting other entered fields.
- Disable duplicate submits while a save request is in progress.
- Preserve user-entered data on validation or server failure.
- Close the form only after a successful save or an explicit discard decision.
- On success, show a short confirmation and update the transaction list and financial summaries from the confirmed response.
- On failure, show the cause in plain Indonesian and leave the form editable.
- If the form has unsaved data, closing it must provide a clear discard/continue choice.

AI categorization:

- Display a category suggestion only when the model returns a valid prediction.
- Make it clear when a category was suggested automatically.
- Allow a one-tap correction.
- If confidence is below the configured threshold, ask the user to choose rather than silently committing the suggestion.
- The form must remain fully usable when the ML service is unavailable.

### 6.4 Budgets — Anggaran

**Primary purpose:** Show how much has been spent, what remains, and which categories need attention.

Content order:

1. Screen title and short purpose statement.
2. Overall budget amount and reporting period.
3. Used amount, remaining amount, and percentage used.
4. Per-category budget rows with used/limit values and progress indicators.
5. A concise list of categories that need attention.
6. An action to create or edit a budget.

Behavior:

- Always display used and remaining values alongside a progress bar.
- Progress must be capped visually at 100% while still showing the actual overspent amount and percentage in text.
- Use status labels such as “Dalam batas,” “Mendekati batas,” or “Melebihi batas.”
- Warning thresholds must be consistent and configurable rather than chosen differently in each screen.
- Opening an alert should take the user to the affected budget or its relevant transactions.
- Internal account transfers must not count as external expenses.

### 6.5 Insights — Insight

**Primary purpose:** Explain meaningful financial patterns and help users decide what to inspect next.

Content order:

1. Screen title and a plain-language description of the feature.
2. One leading insight with period and evidence.
3. Additional findings, ranked by relevance and actionability.
4. Forecast information, when supported by sufficient data.
5. Links to the relevant budget, report, or transaction evidence.

Every insight card should state:

- What was detected.
- Which time period it covers.
- The numeric evidence or comparison behind the statement.
- Whether the value is a historical fact, a statistical detection, or a forecast.
- The next useful action, if any.

AI/ML presentation rules:

- Historical totals must come from deterministic backend calculations, not LLM-generated arithmetic.
- A forecast must be labelled “Perkiraan” or “Estimasi,” not presented as a certain outcome.
- Show uncertainty ranges only if the model can produce and validate them.
- Do not show an arbitrary confidence percentage solely to make AI look sophisticated.
- Explain warnings using evidence, such as “87% anggaran terpakai; sisa Rp88.000,” rather than vague statements like “AI menemukan sesuatu.”
- Do not label unusual spending as fraud. Use wording such as “Pengeluaran ini berbeda dari pola biasanya; periksa transaksi.”
- If the history is insufficient, say so and explain what period or amount of data is missing.
- A user must be able to open the transaction or aggregate evidence supporting an insight.
- If there is no meaningful insight, use a calm empty state rather than generating filler.

### 6.6 Profile and settings — Profil

**Primary purpose:** Give users control over account preferences, privacy, and their data.

Include:

- User identity and account management.
- Light/dark/system theme choice.
- Privacy and security settings.
- Export data.
- Data deletion and account deletion, when supported by the backend.
- App version and legal/privacy information.

Behavior:

- Do not expose full account or transaction details on the profile screen without a clear need.
- Explain what export or deletion does before the action is confirmed.
- Show confirmation and failure states for actions that affect user data.
- Theme choice should persist across app launches when preferences storage is implemented.

## 7. Shared components

### 7.1 App bar

- Brand mark and wordmark on primary screens.
- Avatar button navigates to Profile.
- Icon-only actions must have semantic labels and touch targets meeting the accessibility baseline.
- Avoid placing low-frequency controls in the top bar if they compete with the screen’s primary task.

### 7.2 Financial summary surface

- Label the metric and its period.
- Use consistent currency formatting and numeral alignment.
- Keep the primary number dominant; secondary change/comparison is smaller but still legible.
- Do not use color to imply a trend unless the direction and comparison period are stated.

### 7.3 Transaction row

- Use a minimum practical row height with sufficient touch area.
- Align all monetary values to a consistent end edge.
- Make title, category, and date semantically distinct.
- Do not use merchant logos as required dependencies for readability.

### 7.4 Status chip

- Short text plus color/surface treatment.
- Avoid more than one status chip in a compact row unless each represents distinct information.
- Status wording must be consistent across dashboard, budgets, and insights.

### 7.5 Progress bar

- Always pair with a text value or accessible semantic value.
- Expose percentage and label to screen readers.
- Handle values over the limit and zero-limit budgets safely.
- Avoid gradients or ornamental chart effects.

### 7.6 Bottom sheet and modal

- Respect system bars and bottom safe area.
- Allow scrolling when the keyboard is open.
- Place the primary action where it remains discoverable above the keyboard where feasible.
- Avoid nested sheets and dialog chains for simple tasks.

## 8. The 10 usability heuristics applied to MonTra

This section applies Nielsen’s 10 usability heuristics to MonTra. The source defines these as broad evaluation principles, not a substitute for usability testing.

### 1. Visibility of system status

Users should understand whether an action is in progress, completed, or failed.

MonTra implementation:

- Show saving progress for transaction submission.
- Confirm successful saves and clearly present validation/server errors.
- Mark forecasts with the relevant period and freshness where meaningful.
- Indicate the selected tab, active filter, and chosen reporting period.
- Distinguish cached or stale information from freshly loaded information when the distinction matters.

### 2. Match between the system and the real world

Use familiar personal-finance language and formats.

MonTra implementation:

- Use Indonesian labels such as “Pemasukan,” “Pengeluaran,” “Sisa anggaran,” and “Perkiraan akhir bulan.”
- Format dates and IDR values for the Indonesian locale.
- Show real dates and readable category names rather than internal identifiers.
- Use financial terms consistently; explain unavoidable technical terms in plain language.

### 3. User control and freedom

Users should be able to correct mistakes and leave unwanted flows.

MonTra implementation:

- Allow users to edit transaction details and change auto-suggested categories.
- Provide clear cancel/close behavior for forms and sheets.
- Warn before discarding unsaved data.
- Provide undo for deletions when practical.
- Keep privacy controls such as hiding balances under user control.

### 4. Consistency and standards

Similar things should work and look similar across the app.

MonTra implementation:

- Use one consistent placement for amount, category, and date across transaction lists.
- Reuse budget statuses, colors, labels, spacing, and chart conventions.
- Use Android-standard back behavior, keyboard behavior, touch interactions, and semantic controls.
- Use the shared design tokens rather than one-off widget styling.

### 5. Error prevention

Prevent problems before users commit them.

MonTra implementation:

- Validate positive amounts and required fields before submission.
- Prevent duplicate form submission.
- Check that a category is valid for the selected transaction type.
- Confirm destructive actions and unsaved-data dismissal.
- Do not allow an internal transfer to be accidentally counted as ordinary income and expense simultaneously.

### 6. Recognition rather than recall

Make choices visible when users need them.

MonTra implementation:

- Display category suggestions and recent accounts in context.
- Show active filters, reporting period, budget limit, used amount, and remaining amount.
- Use readable transaction labels and visible category names.
- Remember reasonable preferences without hiding how they were set.

### 7. Flexibility and efficiency of use

Support common tasks quickly without confusing first-time users.

MonTra implementation:

- Make Add transaction available from every primary tab.
- Keep transaction entry short and put optional fields behind a secondary action where appropriate.
- Provide search and simple filters.
- Use sensible defaults while keeping every default editable.
- Consider recurring transactions and import tools as later efficiency features, not MVP blockers.

### 8. Aesthetic and minimalist design

Every visual element should help the task.

MonTra implementation:

- Keep the home screen focused on balance, cash flow, budget status, and recent activity.
- Show the most actionable insight first; place supporting analysis in the Insight screen.
- Remove decorative AI graphics, redundant cards, unnecessary badges, and unlabelled chart decoration.
- Use whitespace and alignment to create hierarchy instead of adding more visual effects.

### 9. Help users recognize, diagnose, and recover from errors

Error messages must explain what happened and what the user can do next.

MonTra implementation:

- Replace vague errors with direct Indonesian messages, such as “Transaksi belum tersimpan. Periksa koneksi lalu coba lagi.”
- Preserve form contents when saving fails.
- Explain invalid values next to the relevant field.
- Provide a retry action for transient network/service failures.
- If AI inference fails, tell the user that automatic suggestions are unavailable and allow manual categorization.

### 10. Help and documentation

Provide help when it is needed, without forcing users to read a manual to do simple tasks.

MonTra implementation:

- Add concise helper text for uncertain concepts, such as the meaning of a forecast or how budget usage is calculated.
- Explain that the recorded balance depends on transactions and opening balances entered or imported by the user.
- Provide a privacy explanation for what transaction data is used for AI features.
- Use contextual help on relevant screens and maintain a compact FAQ/support destination as the product grows.

## 9. Android accessibility and responsive behavior

### 9.1 Touch and interaction targets

- Aim for at least **48 × 48 dp** for interactive touch targets, with enough separation to prevent accidental taps. This is consistent with Android accessibility guidance.
- Small icons may have a smaller visible glyph but must retain a sufficiently large touch region.
- Do not place closely spaced controls where their touch areas overlap.
- Ensure tappable list rows have a clear interactive area.

### 9.2 Text and contrast

- Use `sp` for text sizes and `dp` for layout dimensions.
- Test larger Android font settings; avoid fixed heights that clip labels or transaction titles.
- Aim for WCAG AA contrast for normal text (4.5:1) and large text (3:1), including disabled or secondary states where applicable. Do not assume the palette passes without measurement.
- Do not encode status using color alone; pair color with label, icon, sign, or shape.

### 9.3 Screen readers and semantics

- Give every icon-only action a meaningful semantic label.
- Expose balances with currency and context, not as a sequence of isolated digits.
- Give charts a concise textual summary and expose important values in an accessible structure.
- Expose budget progress as a label and numeric value.
- Avoid announcing purely decorative elements.
- Verify TalkBack order follows the intended visual and task sequence.

### 9.4 Screen sizes and system UI

- Support common Android widths and tall screens without horizontal scrolling for core screens.
- Respect status bar, navigation bar, display cutout, and IME/keyboard insets.
- Use scrollable content where content can exceed the viewport.
- Avoid fixed screen heights and absolute positioning for core app layouts.
- Check landscape or split-screen behavior where supported, even if portrait is the primary layout.

### 9.5 Motion

- Provide immediate feedback but do not rely on motion to communicate important state.
- Use reduced or no motion when requested by system settings where possible.
- Avoid transitions that delay transaction entry or block controls.

## 10. Data, privacy, and trustworthy AI UX

Money data is sensitive. The interface must make the user’s control over data visible and understandable.

- Keep balance concealment easy to find.
- Avoid showing complete transaction descriptions or totals in unnecessary notifications.
- Ask for permissions only when the related feature is invoked and explain why.
- Never imply that MonTra is connected to a bank or wallet unless a real, authorized integration exists.
- Clearly identify demo data during prototype review.
- Explain that AI categorization is a suggestion and can be corrected.
- Show the date range for historical comparisons and forecasts.
- Keep calculation evidence accessible from each financial insight.
- Never invent account balances, historical events, spending reasons, or model confidence values.
- Do not send user financial data to an external model unless the processing purpose, data minimization, provider terms, and user notice/consent requirements have been addressed.
- Make the app’s data export, retention, and deletion behavior understandable before confirmation.

## 11. Content and microcopy guidelines

Write in clear, natural Indonesian. Use short sentences and concrete financial language. Avoid generic marketing language inside the app.

Preferred examples:

- “Belanja kebutuhan sudah memakai 87% anggaran.”
- “Sisa anggaran Rp88.000 sampai 31 Oktober.”
- “Pengeluaran bulan ini 6,4% lebih rendah daripada September.”
- “Perkiraan akhir bulan: Rp3.500.000. Hasil dapat berubah saat transaksi baru dicatat.”
- “Kategori belum pasti. Pilih kategori yang sesuai.”
- “Transaksi belum tersimpan. Periksa koneksi dan coba lagi.”

Avoid:

- “AI tahu kamu boros.”
- “Keuanganmu pasti aman.”
- “AI mendeteksi penipuan” without substantiating evidence.
- “Amazing!”, “Level up!”, or other generic gamification copy that adds no useful information.

For dates and comparisons, state the timeframe clearly. For amounts, use one consistent IDR format throughout the app.

## 12. Screen states and edge cases

Every data-driven screen must define more than its happy path.

### Loading

- Use stable placeholders or a compact progress indicator.
- Avoid flashing empty-state messages while a request is still loading.
- Do not animate skeletons excessively.

### Empty data

- Explain why the screen is empty and give one relevant next action.
- Example on Transactions: “Belum ada transaksi untuk periode ini” with “Catat transaksi.”
- Example on Insights: “Insight akan muncul setelah tersedia cukup riwayat transaksi.”
- Do not invent a sample insight and present it as the user’s actual financial data.

### Search returns no results

- State “Transaksi tidak ditemukan.”
- Provide an easy way to clear the query or filters.
- Do not display the generic first-use empty state.

### Network/API failure

- Explain that data could not be loaded or saved.
- Provide a retry action where safe.
- Preserve unsaved data.
- Do not update financial totals optimistically unless rollback and consistency are fully handled.

### ML unavailable

- Keep core flows working.
- Fall back to manual category selection or baseline logic if appropriate.
- Label the missing AI result accurately; do not use a spinner indefinitely.

### Over-budget and zero-budget cases

- Display the actual over-budget amount and avoid progress bars that overflow their container.
- Define behavior for a zero or archived budget instead of dividing by zero.

### Long text and large numbers

- Support long merchant names, custom categories, large balances, and large text settings without overlapping controls.
- Wrap or truncate with a way to see the full value.
- Keep amount alignment consistent.

## 13. Flutter implementation guidance

The production implementation should translate the design system to Flutter components rather than embedding the HTML prototype inside a WebView.

Recommended implementation rules:

- Centralize the palette, typography, spacing, shape, and elevation tokens.
- Use `ThemeData` and Material components where they fit; customize them through theme tokens rather than repeated ad hoc styles.
- Use `SafeArea`, keyboard-aware scrollable forms, and appropriate insets.
- Keep widgets small and single-purpose: app bar, balance surface, flow summary, transaction row, budget progress row, insight card, and bottom navigation.
- Separate presentation from data access and business rules.
- Keep money formatting in one tested utility; do not format currencies differently in each widget.
- Use semantic labels and minimum touch dimensions for custom controls.
- Handle asynchronous states explicitly: initial, loading, success, empty, and error.
- Use one source of truth for transaction data, with state invalidation after successful mutations.
- Do not calculate authoritative account balances or forecast totals in presentation widgets.
- Avoid hardcoded sample figures in production widgets. Put demo fixtures in clearly named mock repositories or test data.
- Use feature-level folders or another consistent organization; avoid scattering UI, API requests, and domain calculations in a single large file.

Suggested UI structure:

```text
apps/mobile/lib/
├── app/
│   ├── app.dart
│   ├── router.dart
│   └── theme/
│       ├── app_theme.dart
│       ├── app_colors.dart
│       ├── app_typography.dart
│       ├── app_spacing.dart
│       └── app_shapes.dart
├── core/
│   ├── formatting/
│   ├── networking/
│   ├── accessibility/
│   └── widgets/
└── features/
    ├── home/
    ├── transactions/
    ├── budgets/
    ├── insights/
    └── profile/
```

This is a suggested organization, not a requirement to create empty abstractions. Add modules when they contain real responsibilities.

## 14. Quality assurance and design review

Review all core flows in both light and dark themes on a small Android viewport and a larger device. Test the actual Flutter build, not only a desktop browser preview.

### Required usability tests

1. A first-time user can identify how to add an expense without assistance.
2. A user can record an expense and recognize whether it was saved.
3. A user can find a transaction using search and filters.
4. A user can determine the remaining budget for a category.
5. A user can tell the difference between historical spending and a forecast.
6. A user can explain the evidence for a displayed insight.
7. A user can hide and restore the displayed balance.
8. A user can recover from an invalid amount or failed save without re-entering all data.

### Visual QA checklist

- [ ] Palette matches the defined tokens in light and dark modes.
- [ ] Hierarchy is visible without excessive shadows or decorative cards.
- [ ] Money values use consistent IDR formatting.
- [ ] Text remains readable with increased Android font scale.
- [ ] Touch targets meet the 48 × 48 dp baseline where applicable.
- [ ] TalkBack labels, focus order, and chart descriptions are useful.
- [ ] Keyboard and system bars do not cover form fields or primary actions.
- [ ] Empty, loading, error, over-budget, and ML-unavailable states exist.
- [ ] Forecasts and AI insights clearly indicate period, evidence, and uncertainty.
- [ ] No production screen contains fake transactions or unsupported AI claims.
- [ ] Core transaction flows work when ML features are disabled.
- [ ] No horizontal overflow or clipped large monetary values appear on supported device sizes.

### Acceptance criteria

The design is implementation-ready when all screen structures and interaction rules are mapped to Flutter components, the design tokens are centralized, required edge states are specified, and the core usability tests can be completed without facilitator explanation. Visual polish alone is not sufficient for acceptance.

## 15. References

1. Nielsen Norman Group, [10 Usability Heuristics for User Interface Design](https://www.nngroup.com/articles/ten-usability-heuristics/). The 10 heuristics are used as the usability review framework in this document.
2. Android Developers, [Make apps more accessible](https://developer.android.com/guide/topics/ui/accessibility/apps). Guidance for Android touch targets, semantics, and accessibility.
3. Android Accessibility Help, [Touch target size](https://support.google.com/accessibility/android/answer/7101858). Recommends touch targets of at least 48 × 48 dp.

---

## 16. Design principle summary

MonTra’s Android interface should be judged by whether people can record, understand, and verify their financial activity with minimal friction. The chosen style—warm neutral surfaces, forest-green emphasis, restrained lime accents, readable numbers, and evidence-led insights—supports that goal. Any design decision that adds visual novelty while reducing clarity, accessibility, control, or trust should be rejected.
