# Store Customization — Critères d'acceptation (MVP)

## User Story
En tant que commerçant, je veux personnaliser l’apparence de ma boutique (logo, couleurs, thème), afin d’adapter la boutique à l’identité de ma marque.

## Critères d’acceptation obligatoires

1. **Branding éditable**
   - Le commerçant peut modifier: `logo`, `theme_primary_color`, `theme_secondary_color`, `theme_mode`, `theme_font_family`.
   - Les valeurs sont persistées au niveau du store.

2. **Builder drag & drop**
   - Le commerçant peut réordonner les sections de vitrine via glisser-déposer.
   - Le commerçant peut activer/désactiver chaque section.
   - Le layout est persisté dans un brouillon (`storefront_layout_draft`).

3. **Workflow draft/publish**
   - Le commerçant peut sauvegarder un brouillon sans publication.
   - Le commerçant peut publier la personnalisation.
   - La publication copie le layout brouillon vers la version publiée (`storefront_layout_published`) et incrémente `theme_version`.

4. **Preview fonctionnelle**
   - Une preview de la vitrine reflète le brouillon courant.
   - Les couleurs et le mode de thème sont visibles sur la preview.

5. **Multi-tenant / sécurité**
   - Les endpoints restent sous authentification tenant.
   - La personnalisation est scoped au store du tenant.

## Endpoints MVP
- `PUT /api/stores/:id` (save draft + branding)
- `POST /api/stores/:id/customization/publish` (publish)
- `GET /api/stores/:id` (read settings)

## Remarques
- Ce MVP couvre la personnalisation dashboard + publication versionnée.
- L’exposition publique cross-tenant par slug n’est pas incluse dans ce lot (architecture actuelle en schémas tenant).