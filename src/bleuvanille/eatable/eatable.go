/*

Problématique

la personnalisation
collecte de données
gestion des tags :  cuit, mangé froid, à cuire avec autocuiseur, ...
detection automatique de paramètre : quantité par personne, type de cuisine, difficulté, detection auto d'ingrédient manquand ou optionnel
liste de liens

Fonction

creation d'un ingredient
mise à jour du père d'un ingrédient
modification d'un ou plusieurs champs d'un ingredient
changement de status d'un ingrédient
récupération d'un ingrédient
récupération des enfants d'un ingrédient
récupération des ingrédients proches d'un ingrédient


création

Ingredient

nom
status
type
date de creation
date de modification
description


Collection children;
Ingredient parent;

Nutrient nutrient;

preferences;
PriceBean price;

List proportions;
Double indispensable = Quantity.NONSIGNIFICANT_QTY;
Double eatableportion = 1.0;

Date seasonbegin;
Date seasonend;

Integer unit = Quantity.TYPE_UNIT_WEIGHT;
String boughtunit = Quantity.TYPE_UNIT_WEIGHT_STR;
double weightPerUnit = 0.1;
double density = 1.0;

String synonym;

long peremption = 0;
int frequency = 0;
Integer vege = Ingredient.TYPE_INGREDIENT_VEGE_NO;

RECETTE
	
Date rest;
Date preparation;
Date cooking;

Integer level = Recipe.LEVEL_EASY;
Integer occasion = Recipe.OCCASION_NORMAL;

Cooker owner;

Quantity quantity;
List<Proportion> proportions;
List<Realization> realizations;
Set<RecipeCategory> categories;


String origin;

Integer boost;

boolean presurecooking = false;
Boolean coldeatingIntern = false;
Boolean heatableIntern = false;

*/
package eatable
