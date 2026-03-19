'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { useAuth } from '@/lib/hooks/use-auth';
import { useCreateProduct, useCategories } from '@/lib/hooks/use-api';

const productSchema = z.object({
  title: z.string().min(1, 'Title is required'),
  slug: z.string().optional(),
  description: z.string().optional(),
  price: z.number().min(0, 'Price must be >= 0'),
  sale_price: z.number().optional(),
  status: z.enum(['draft', 'published', 'archived']),
  visibility: z.enum(['public', 'private']),
  track_stock: z.boolean(),
  stock: z.number().min(0, 'Stock must be >= 0'),
  currency: z.string().length(3, 'Currency must be 3 characters'),
  sku: z.string().optional(),
  weight: z.number().optional(),
  dimensions: z.string().optional(),
  brand: z.string().optional(),
  tax_class: z.string().optional(),
  category_id: z.string().optional(),
});

type ProductForm = z.infer<typeof productSchema>;

export default function NewProductPage() {
  const router = useRouter();
  const { currentStore } = useAuth();
  const storeId = currentStore?.id || '';
  const [error, setError] = useState<string>('');

  const createProductMutation = useCreateProduct();
  const { data: categories } = useCategories(storeId);

  const { register, handleSubmit, formState: { errors } } = useForm<ProductForm>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      status: 'draft',
      visibility: 'public',
      track_stock: true,
      stock: 0,
      currency: 'USD',
    },
  });

  const onSubmit = async (data: ProductForm) => {
    if (!storeId) {
      setError('Select a store before creating products.');
      return;
    }

    try {
      setError('');
      await createProductMutation.mutateAsync({ storeId, data });
      router.push('/dashboard/products');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create product');
    }
  };


  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-3xl">
        <CardHeader>
          <CardTitle>Create Product</CardTitle>
          <CardDescription>
            Add a new product to your store.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="title">Product Title *</Label>
                <Input id="title" placeholder="Cool T-Shirt" {...register('title')} />
                {errors.title && <p className="text-sm text-red-600">{errors.title.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="slug">Slug</Label>
                <Input id="slug" placeholder="cool-t-shirt" {...register('slug')} />
                {errors.slug && <p className="text-sm text-red-600">{errors.slug.message}</p>}
              </div>

              <div className="space-y-2 lg:col-span-2">
                <Label htmlFor="description">Description</Label>
                <Input id="description" placeholder="Product description" {...register('description')} />
                {errors.description && <p className="text-sm text-red-600">{errors.description.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="price">Price *</Label>
                <Input id="price" type="number" step="0.01" placeholder="0.00" {...register('price', { valueAsNumber: true })} />
                {errors.price && <p className="text-sm text-red-600">{errors.price.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="sale_price">Sale Price</Label>
                <Input id="sale_price" type="number" step="0.01" placeholder="0.00" {...register('sale_price', { valueAsNumber: true })} />
                {errors.sale_price && <p className="text-sm text-red-600">{errors.sale_price.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="currency">Currency *</Label>
                <Input id="currency" placeholder="USD" {...register('currency')} />
                {errors.currency && <p className="text-sm text-red-600">{errors.currency.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="stock">Stock *</Label>
                <Input id="stock" type="number" placeholder="0" {...register('stock', { valueAsNumber: true })} />
                {errors.stock && <p className="text-sm text-red-600">{errors.stock.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="status">Status *</Label>
                <select id="status" {...register('status')} className="flex h-9 rounded-md border border-input bg-transparent px-3 py-1 text-site-text shadow-sm transition-colors file:border-0 file:bg-transparent file:text-site-text file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50">
                  <option value="draft">Draft</option>
                  <option value="published">Published</option>
                  <option value="archived">Archived</option>
                </select>
                {errors.status && <p className="text-sm text-red-600">{errors.status.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="visibility">Visibility *</Label>
                <select id="visibility" {...register('visibility')} className="flex h-9 rounded-md border border-input bg-transparent px-3 py-1 text-site-text shadow-sm transition-colors file:border-0 file:bg-transparent file:text-site-text file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50">
                  <option value="public">Public</option>
                  <option value="private">Private</option>
                </select>
                {errors.visibility && <p className="text-sm text-red-600">{errors.visibility.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="sku">SKU</Label>
                <Input id="sku" placeholder="SKU-001" {...register('sku')} />
                {errors.sku && <p className="text-sm text-red-600">{errors.sku.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="weight">Weight</Label>
                <Input id="weight" type="number" step="0.01" placeholder="1.5" {...register('weight', { valueAsNumber: true })} />
                {errors.weight && <p className="text-sm text-red-600">{errors.weight.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="brand">Brand</Label>
                <Input id="brand" placeholder="Brand name" {...register('brand')} />
                {errors.brand && <p className="text-sm text-red-600">{errors.brand.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="dimensions">Dimensions</Label>
                <Input id="dimensions" placeholder="10x10x10" {...register('dimensions')} />
                {errors.dimensions && <p className="text-sm text-red-600">{errors.dimensions.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="tax_class">Tax Class</Label>
                <Input id="tax_class" placeholder="standard" {...register('tax_class')} />
                {errors.tax_class && <p className="text-sm text-red-600">{errors.tax_class.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="category_id">Category (optional)</Label>
                <select id="category_id" {...register('category_id')} className="flex h-9 rounded-md border border-input bg-transparent px-3 py-1 text-site-text shadow-sm transition-colors file:border-0 file:bg-transparent file:text-site-text file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50">
                  <option value="">Select a category</option>
                  {categories?.map((cat) => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
                {errors.category_id && <p className="text-sm text-red-600">{errors.category_id.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="track_stock" className="flex items-center gap-2">
                  <input type="checkbox" id="track_stock" {...register('track_stock')} className="h-4 w-4" />
                  <span>Track Stock</span>
                </Label>
              </div>
            </div>

            <Button type="submit" className="w-full" disabled={createProductMutation.isPending}>
              Create Product
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
