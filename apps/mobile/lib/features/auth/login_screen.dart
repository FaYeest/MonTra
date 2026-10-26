import 'package:flutter/material.dart';

import '../../app/theme/app_colors.dart';
import '../../app/theme/app_shapes.dart';
import '../../app/theme/app_spacing.dart';
import '../../core/widgets/montra_scope.dart';
import 'register_screen.dart';

/// Sign-in screen for MonTra.
class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    if (!_formKey.currentState!.validate()) return;

    final auth = MontraScope.authOf(context);
    final ledger = MontraScope.ledgerOf(context);

    final success = await auth.login(
      _emailController.text,
      _passwordController.text,
    );

    if (success && mounted) {
      await ledger.refreshAll();
    }
  }

  @override
  Widget build(BuildContext context) {
    final tokens = Theme.of(context).extension<MontraColorTokens>()!;
    final textTheme = Theme.of(context).textTheme;
    final auth = MontraScope.authOf(context);

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(
              horizontal: MontraSpacing.pageHorizontal,
              vertical: MontraSpacing.xl,
            ),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: <Widget>[
                  // Logo / App header
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: <Widget>[
                      Container(
                        padding: const EdgeInsets.all(MontraSpacing.sm),
                        decoration: BoxDecoration(
                          color: tokens.primary.withValues(alpha: 0.12),
                          borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
                        ),
                        child: Icon(
                          Icons.account_balance_wallet_rounded,
                          color: tokens.primary,
                          size: 32,
                        ),
                      ),
                      const SizedBox(width: MontraSpacing.sm),
                      Text(
                        'MonTra',
                        style: textTheme.headlineMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: tokens.primary,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: MontraSpacing.xl),
                  Text(
                    'Masuk ke Akun Anda',
                    style: textTheme.headlineSmall,
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: MontraSpacing.xxs),
                  Text(
                    'Catat dan pantau transaksi finansial dengan transparan.',
                    style: textTheme.bodyMedium?.copyWith(
                      color: tokens.textSecondary,
                    ),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: MontraSpacing.xl),

                  // Error banner if any
                  if (auth.errorMessage != null) ...[
                    Container(
                      padding: const EdgeInsets.all(MontraSpacing.md),
                      decoration: BoxDecoration(
                        color: tokens.expense.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(MontraShapes.radiusSm),
                        border: Border.all(
                          color: tokens.expense.withValues(alpha: 0.3),
                          width: MontraShapes.outlineWidth,
                        ),
                      ),
                      child: Row(
                        children: <Widget>[
                          Icon(Icons.error_outline, color: tokens.expense, size: 20),
                          const SizedBox(width: MontraSpacing.xs),
                          Expanded(
                            child: Text(
                              auth.errorMessage!,
                              style: textTheme.bodySmall?.copyWith(
                                color: tokens.expense,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: MontraSpacing.md),
                  ],

                  // Email field
                  TextFormField(
                    controller: _emailController,
                    keyboardType: TextInputType.emailAddress,
                    textInputAction: TextInputAction.next,
                    decoration: const InputDecoration(
                      labelText: 'Email',
                      hintText: 'nama@domain.com',
                      prefixIcon: Icon(Icons.email_outlined),
                    ),
                    validator: (val) {
                      if (val == null || val.trim().isEmpty) {
                        return 'Email wajib diisi';
                      }
                      if (!val.contains('@') || !val.contains('.')) {
                        return 'Format email tidak valid';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: MontraSpacing.md),

                  // Password field
                  TextFormField(
                    controller: _passwordController,
                    obscureText: _obscurePassword,
                    textInputAction: TextInputAction.done,
                    onFieldSubmitted: (_) => _handleLogin(),
                    decoration: InputDecoration(
                      labelText: 'Kata Sandi',
                      prefixIcon: const Icon(Icons.lock_outline),
                      suffixIcon: IconButton(
                        icon: Icon(
                          _obscurePassword
                              ? Icons.visibility_outlined
                              : Icons.visibility_off_outlined,
                        ),
                        onPressed: () {
                          setState(() {
                            _obscurePassword = !_obscurePassword;
                          });
                        },
                      ),
                    ),
                    validator: (val) {
                      if (val == null || val.isEmpty) {
                        return 'Kata sandi wajib diisi';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: MontraSpacing.lg),

                  // Submit button
                  SizedBox(
                    height: MontraSpacing.minTouchTarget,
                    child: ElevatedButton(
                      onPressed: auth.isLoading ? null : _handleLogin,
                      child: auth.isLoading
                          ? const SizedBox(
                              width: 24,
                              height: 24,
                              child: CircularProgressIndicator(
                                strokeWidth: 2.5,
                                color: Colors.white,
                              ),
                            )
                          : const Text('Masuk'),
                    ),
                  ),
                  const SizedBox(height: MontraSpacing.lg),

                  // Switch to register
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: <Widget>[
                      Text(
                        'Belum punya akun?',
                        style: textTheme.bodyMedium?.copyWith(
                          color: tokens.textSecondary,
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          auth.clearError();
                          Navigator.of(context).push(
                            MaterialPageRoute<void>(
                              builder: (_) => const RegisterScreen(),
                            ),
                          );
                        },
                        child: const Text('Daftar di sini'),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
