# Recipe Knowledge Base Example

This example configures the MCP server for managing culinary recipes, ingredients, techniques, and meal planning.

## Use Case

Build a queryable knowledge graph of culinary information:
- Recipe dependencies (ingredients, base recipes, techniques)
- Meal planning and pairing suggestions
- Recipe variations and substitutions
- Dietary restrictions and allergen tracking
- Cuisine types and cooking methods

## Configuration

### MCP Tools Generated

With this configuration, the server exposes:
- `add_recipe` - Create a new recipe with ingredients and pairings
- `get_recipe` - Retrieve a recipe with its full relationship graph
- `list_recipes` - Browse recipes by tags (cuisine, dietary, etc.)
- `update_recipe` - Modify an existing recipe
- `delete_recipe` - Remove a recipe

### Relationship Types

**requires** (backward direction)
- Points to ingredients, base recipes, or techniques needed
- Example: "chicken-parmesan" requires "marinara-sauce" and "breaded-chicken"
- Can reference both ingredient nodes and sub-recipe nodes

**produces** (forward direction)
- Points to final dishes or meal components this recipe creates
- Example: "roast-chicken" produces "chicken-dinner" and "chicken-stock"
- Useful for tracking yields and meal planning

**pairs_with** (no temporal direction)
- Recipes that complement this one in a meal
- Example: "grilled-salmon" pairs_with "asparagus" and "lemon-rice"
- Bidirectional relationships for meal planning

**variations** (no temporal direction)
- Alternative versions or preparation methods
- Example: "chocolate-chip-cookies" variations include "gluten-free-chocolate-chip-cookies"
- Helps with dietary substitutions

## Example Recipe

```yaml
id: chicken-parmesan
name: Chicken Parmesan
summary: Classic Italian-American breaded chicken with marinara and mozzarella
description: |
  Crispy breaded chicken cutlets topped with marinara sauce and melted
  mozzarella cheese, baked until golden and bubbly.

  Prep time: 20 minutes
  Cook time: 25 minutes
  Serves: 4

  Instructions:
  1. Bread chicken cutlets with panko and parmesan
  2. Pan-fry until golden brown
  3. Top with marinara sauce and mozzarella
  4. Bake at 400°F for 15 minutes until cheese melts
  5. Garnish with fresh basil

  Nutrition: 450 cal, 28g protein, 18g fat, 35g carbs
tags:
  - italian
  - main-course
  - chicken
  - baked
  - comfort-food
edges:
  requires:
    - chicken-breast
    - panko-breadcrumbs
    - parmesan-cheese
    - mozzarella-cheese
    - marinara-sauce
  produces:
    - chicken-parmesan-dinner
  pairs_with:
    - garlic-bread
    - caesar-salad
    - spaghetti-marinara
  variations:
    - eggplant-parmesan
    - veal-parmesan
    - gluten-free-chicken-parmesan
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-01-15T10:30:00Z
```

## Example Ingredient Node

```yaml
id: marinara-sauce
name: Marinara Sauce
summary: Classic Italian tomato sauce
description: |
  Simple Italian tomato sauce with garlic, basil, and olive oil.

  Prep time: 5 minutes
  Cook time: 30 minutes
  Yield: 4 cups

  Instructions:
  1. Sauté garlic in olive oil
  2. Add crushed tomatoes and simmer
  3. Season with basil, salt, and pepper
  4. Simmer 30 minutes until thickened
tags:
  - sauce
  - italian
  - tomato-based
  - vegetarian
  - vegan
edges:
  requires:
    - crushed-tomatoes
    - garlic
    - olive-oil
    - basil
    - salt
    - black-pepper
  produces:
    - marinara-sauce
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-01-15T10:30:00Z
```

## Common Tags

Suggested tags for organizing recipes:

- **Cuisine**: `italian`, `french`, `mexican`, `thai`, `chinese`, `indian`, `american`
- **Course**: `appetizer`, `main-course`, `side-dish`, `dessert`, `beverage`, `sauce`
- **Dietary**: `vegetarian`, `vegan`, `gluten-free`, `dairy-free`, `keto`, `paleo`
- **Protein**: `chicken`, `beef`, `pork`, `fish`, `seafood`, `tofu`, `legumes`
- **Method**: `baked`, `grilled`, `fried`, `slow-cooked`, `pressure-cooked`, `raw`
- **Meal**: `breakfast`, `lunch`, `dinner`, `snack`
- **Time**: `quick` (< 30 min), `moderate` (30-60 min), `involved` (> 60 min)
- **Occasion**: `weeknight`, `holiday`, `party`, `meal-prep`, `comfort-food`

## Recipe Graph Example

A meal planning graph might look like:

```
garlic ──┐
olive-oil ┼──> marinara-sauce ──┐
tomatoes ┘                       │
                                 ├──> chicken-parmesan ──pairs_with──> caesar-salad
chicken-breast ──┐               │
panko-breadcrumbs ┼──────────────┘
parmesan-cheese ──┘
```

This structure allows AI assistants to:
- Generate shopping lists from recipes
- Suggest meal pairings
- Find recipes based on available ingredients
- Recommend substitutions for dietary restrictions
- Plan weekly menus with variety

## Running This Example

```bash
# Start the server with this configuration
mcp serve --directory ./docs/examples/recipes

# Or copy to your data directory
cp docs/examples/recipes/mcp.yaml ./my-recipes/
cp docs/examples/recipes/relationships.yaml ./my-recipes/
mcp serve --directory ./my-recipes
```

## Integration Ideas

This recipe graph can be integrated with:
- Meal planning AI assistants
- Smart kitchen devices
- Grocery shopping list generators
- Nutritional analysis tools
- Cooking instruction guides
- Recipe recommendation systems
