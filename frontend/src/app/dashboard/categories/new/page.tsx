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
import { useCreateCategory } from '@/lib/hooks/use-api';
import { useAuth } from '@/lib/hooks/use-auth';

const categorySchema = z.object({
  name: z.string().min(1, 'Name is required'),
  slug: z.string().optional(),
  description: z.string().optional(),
  visibility: z.enum(['public', 'private']),
});

type CategoryForm = z.infer<typeof categorySchema>;

export default function NewCategoryPage() {
  const router = useRouter();
  const { currentStore } = useAuth();
  const storeId = currentStore?.id || '';
  const [error, setError] = useState<string>('');
  const createCategoryMutation = useCreateCategory();

  const { register, handleSubmit, formState: { errors } } = useForm<CategoryForm>({
    resolver: zodResolver(categorySchema),
  });

  const onSubmit = async (data: CategoryForm) => {
    if (!storeId) {
      setError('Select a store first.');
      return;
    }

    try {
      setError('');
      await createCategoryMutation.mutateAsync({ storeId, data });
      router.push('/dashboard/categories');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create category');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <Card className="w-full max-w-lg">
        <CardHeader>
          <CardTitle>Create category</CardTitle>
          <CardDescription>
            Add a category to organize products.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="space-y-2">
              <Label htmlFor="name">Name *</Label>
              <Input id="name" placeholder="T-Shirts" {...register('name')} />
              {errors.name && <p className="text-sm text-red-600">{errors.name.message}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="slug">Slug</Label>
              <Input id="slug" placeholder="t-shirts" {...register('slug')} />
              {errors.slug && <p className="text-sm text-red-600">{errors.slug.message}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Input id="description" placeholder="A descriptive summary" {...register('description')} />
              {errors.description && <p className="text-sm text-red-600">{errors.description.message}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="visibility">Visibility *</Label>
              <select id="visibility" {...register('visibility')} className="flex h-9 rounded-md border border-input bg-transparent px-3 py-1 text-site-text shadow-sm transition-colors file:border-0 file:bg-transparent file:text-site-text file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50">
                <option value="public">Public</option>
                <option value="private">Private</option>
              </select>
              {errors.visibility && <p className="text-sm text-red-600">{errors.visibility.message}</p>}
            </div>

            <Button type="submit" className="w-full" disabled={createCategoryMutation.isPending}>
              Create category
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
